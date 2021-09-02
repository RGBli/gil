package gil

import (
	"strings"
)

// 前缀树结构
type node struct {
	// 待匹配路由
	pattern  string
	// 路由中的一部分
	part     string
	children []*node
	// 匹配 part 是否为模糊匹配，支持 : 和 * 两种模糊匹配方式
	isWild   bool
}

// 向前缀树递归插入节点，parts 数组是从待匹配路由 pattern 中分割出来的
func (n *node) insert(pattern string, parts []string, height int) {
	// 递归出口，仅在目标节点设置 pattern 属性，中间节点的 pattern 都是 ""
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 搜索适配的节点
func (n *node) search(parts []string, height int) *node {
	// 递归出口
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 如果是中间节点
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

// dfs 遍历前缀树
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// 找到第一个匹配的孩子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 找到所有匹配的孩子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
