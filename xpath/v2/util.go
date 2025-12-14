package v2

import "strings"

// parseSegment 解析路径段类型
func parseSegment(segment string) (NodeType, string) {
	if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
		// 参数节点 {param}
		paramName := segment[1 : len(segment)-1]
		return ParamNode, paramName
	} else if segment == "*" {
		// 单段通配符
		return WildcardNode, ""
	} else if segment == "**" {
		// 多段通配符
		return CatchAllNode, ""
	}
	// 静态节点
	return StaticNode, ""
}

