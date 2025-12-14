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

// splitPath 分割路径
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

// joinPath 连接路径段
func joinPath(segments ...string) string {
	result := ""
	for i, seg := range segments {
		if seg == "" {
			continue
		}
		if i == 0 || seg[0] == '/' {
			result += seg
		} else {
			result += "/" + seg
		}
	}
	if result == "" {
		return "/"
	}
	if result[0] != '/' {
		result = "/" + result
	}
	return result
}
