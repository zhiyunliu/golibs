package v2

import "sync"

// NodeType 定义节点类型
type NodeType int

const (
	StaticNode NodeType = iota
	ParamNode
	WildcardNode
	CatchAllNode
)

// TreeNode 路径树节点
type TreeNode struct {
	segment   string      // 路径段
	nodeType  NodeType    // 节点类型
	paramName string      // 参数名（仅参数节点）
	children  []*TreeNode // 子节点
	info      *NodeInfo   // 节点信息
	isLeaf    bool        // 是否为叶子节点
	childLock sync.RWMutex
}

// MatchResult 匹配结果
type MatchResult struct {
	Matched bool
	Info    *NodeInfo
	Params  map[string]string
}