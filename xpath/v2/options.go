package v2

// Option 路径选项类型
type Option func(*NodeInfo)

// NodeInfo 节点信息，存储所有Option信息
type NodeInfo struct {
	Name string
	Desc string
	Meta map[string]any
}

// WithName 名称选项
func WithName(name string) Option {
	return func(info *NodeInfo) {
		info.Name = name
	}
}

// WithDesc 描述选项
func WithDesc(desc string) Option {
	return func(info *NodeInfo) {
		info.Desc = desc
	}
}

// WithMeta 元数据选项
func WithMeta(meta map[string]any) Option {
	return func(info *NodeInfo) {
		if info.Meta == nil {
			info.Meta = make(map[string]any)
		}
		for k, v := range meta {
			info.Meta[k] = v
		}
	}
}