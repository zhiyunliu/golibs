package v2

// Matcher 路径匹配器
type Matcher struct {
	root *TreeNode
}

// NewMatcher 创建新匹配器
func NewMatcher() *Matcher {
	return &Matcher{
		root: &TreeNode{
			segment:  "/",
			nodeType: StaticNode,
			children: make([]*TreeNode, 0),
			info:     &NodeInfo{},
		},
	}
}

// AddPath 添加路径规则
func (m *Matcher) AddPath(pattern string, options ...Option) error {
	if pattern == "" || pattern[0] != '/' {
		return NewMatcherError("pattern must start with '/'", pattern)
	}

	// 解析并收集所有选项
	info := &NodeInfo{}
	for _, opt := range options {
		opt(info)
	}

	segments := splitPath(pattern)
	m.root.insert(segments, 0, info)
	return nil
}

// MustAddPath 添加路径规则（panic版本）
func (m *Matcher) MustAddPath(pattern string, options ...Option) {
	if err := m.AddPath(pattern, options...); err != nil {
		panic(err)
	}
}

// Match 匹配URL
func (m *Matcher) Match(path string) *MatchResult {
	if path == "" || path[0] != '/' {
		return &MatchResult{Matched: false}
	}

	segments := splitPath(path)
	params := make(map[string]string)

	info := m.root.search(segments, 0, params)
	if info != nil {
		return &MatchResult{
			Matched: true,
			Info:    info,
			Params:  params,
		}
	}

	return &MatchResult{Matched: false}
}

// Walk 遍历所有路径
func (m *Matcher) Walk(visit func(pattern string, info *NodeInfo)) {
	m.root.walk("", visit)
}

// walk 遍历节点
func (n *TreeNode) walk(currentPath string, visit func(pattern string, info *NodeInfo)) {
	if n.isLeaf && n.info != nil {
		visit(currentPath, n.info)
	}

	n.childLock.RLock()
	defer n.childLock.RUnlock()

	for _, child := range n.children {
		childPath := joinPath(currentPath, child.segment)
		child.walk(childPath, visit)
	}
}

// FindPath 查找指定路径模式
func (m *Matcher) FindPath(pattern string) (*NodeInfo, bool) {
	segments := splitPath(pattern)
	params := make(map[string]string)

	info := m.root.findExact(segments, 0, params)
	return info, info != nil
}

// findExact 精确查找（不匹配参数和通配符）
func (n *TreeNode) findExact(segments []string, depth int, params map[string]string) *NodeInfo {
	if depth >= len(segments) {
		if n.isLeaf {
			return n.info
		}
		return nil
	}

	segment := segments[depth]

	n.childLock.RLock()
	defer n.childLock.RUnlock()

	for _, child := range n.children {
		// 只匹配相同的segment
		if child.segment == segment {
			return child.findExact(segments, depth+1, params)
		}
	}

	return nil
}

// insert 插入路由节点
func (n *TreeNode) insert(segments []string, depth int, info *NodeInfo) {
	if depth >= len(segments) {
		// 合并节点信息（如果已存在）
		if n.info != nil {
			// 合并Meta
			if info.Meta != nil {
				if n.info.Meta == nil {
					n.info.Meta = make(map[string]any)
				}
				for k, v := range info.Meta {
					n.info.Meta[k] = v
				}
			}
			// 更新名称和描述（如果提供了）
			if info.Name != "" {
				n.info.Name = info.Name
			}
			if info.Desc != "" {
				n.info.Desc = info.Desc
			}
		} else {
			n.info = info
		}
		n.isLeaf = true
		return
	}

	segment := segments[depth]
	nodeType, paramName := parseSegment(segment)

	// 检查是否已存在相同段
	n.childLock.Lock()
	defer n.childLock.Unlock()

	for _, child := range n.children {
		if child.segment == segment {
			child.insert(segments, depth+1, info)
			return
		}
	}

	// 创建新节点
	newNode := &TreeNode{
		segment:   segment,
		nodeType:  nodeType,
		paramName: paramName,
		children:  make([]*TreeNode, 0),
		info:      &NodeInfo{}, // 初始化空info
	}

	n.children = append(n.children, newNode)
	newNode.insert(segments, depth+1, info)
}

// search 搜索匹配的路由
func (n *TreeNode) search(segments []string, depth int, params map[string]string) *NodeInfo {
	// 如果到达末尾，返回当前节点信息（如果是叶子节点）
	if depth >= len(segments) {
		if n.isLeaf {
			return n.info
		}
		return nil
	}

	segment := segments[depth]

	n.childLock.RLock()
	defer n.childLock.RUnlock()

	// 优先级：静态匹配 > 参数匹配 > 通配符匹配 > 多段匹配

	// 1. 尝试静态匹配
	for _, child := range n.children {
		if child.nodeType == StaticNode && child.segment == segment {
			if info := child.search(segments, depth+1, params); info != nil {
				return info
			}
		}
	}

	// 2. 尝试参数匹配
	for _, child := range n.children {
		if child.nodeType == ParamNode {
			params[child.paramName] = segment
			if info := child.search(segments, depth+1, params); info != nil {
				return info
			}
			// 回溯：删除参数
			delete(params, child.paramName)
		}
	}

	// 3. 尝试单段通配符
	for _, child := range n.children {
		if child.nodeType == WildcardNode {
			if info := child.search(segments, depth+1, params); info != nil {
				return info
			}
		}
	}

	// 4. 尝试多段通配符（**）
	for _, child := range n.children {
		if child.nodeType == CatchAllNode && child.isLeaf {
			// ** 匹配剩余所有段，直接返回节点信息
			return child.info
		}
	}

	return nil
}
