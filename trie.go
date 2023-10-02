package tgw

import "strings"

/*
前缀路由的实现
*/

// 节点结构
type node struct {
	pattern  string  // 待匹配路由，例如 /class/:name/info
	part     string  // 路由中的一部分，例如 :name
	children []*node // 子节点，例如 [info, age, address]
	isVague  bool    // 是否模糊匹配，part 含有 : 或 * 时为 true
}

// 第一个匹配成功的节点（用于插入）
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isVague {
			return child
		}
	}
	return nil
}

// 注册插入
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isVague: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 所有匹配成功的节点（用于查找）
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isVague {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 查找匹配路由
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
