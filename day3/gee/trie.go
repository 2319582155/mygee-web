package gee

import "strings"

type node struct {
	pattern  string  // 路由的完整路径，只有在叶子节点才设置，其余节点为“”
	part     string  // 路由的一部分，如：/p/lang中的lang
	children []*node // [p, lang, doc] 即/p/:lang/doc
	isWild   bool    // 是否是模糊匹配，如：:lang，*File
}

func newNode(part string) *node {
	return &node{
		part:     part,
		children: make([]*node, 0),
		isWild:   part[0] == ':' || part[0] == '*',
	}
}

// 找到n的子节点中第一个能匹配part的，用于前缀树的构建
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if strings.EqualFold(part, child.part) || child.isWild {
			return child
		}
	}
	return nil
}

// 找到n的子节点中所有能匹配part的，用于前缀树的搜索
func (n *node) matchChildren(part string) []*node {
	res := make([]*node, 0)
	for _, child := range n.children {
		if strings.EqualFold(part, child.part) || child.isWild {
			res = append(res, child)
		}
	}
	return res
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = newNode(part)
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

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
		ch := child.search(parts, height+1)
		if ch != nil {
			return ch
		}
	}
	return nil
}
