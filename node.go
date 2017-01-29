package envh

type node struct {
	children []*node
	key      string
	value    string
	hasValue bool
}

func newNode() *node {
	return &node{children: []*node{}}
}

func (n *node) findAllNodesByKey(key string, withValue bool) *[]*node {
	results := []*node{}
	nodes := n.children

	for {
		carry := []*node{}

		for _, node := range nodes {
			if node.key == key {
				if withValue && node.hasValue || !withValue {
					results = append(results, node)
				}
			}

			carry = append(carry, node.children...)
		}

		nodes = carry

		if len(carry) == 0 {
			return &results
		}
	}
}

func (n *node) findNodeByKeyChain(keyChain *[]string) (*node, bool) {
	if len(*keyChain) == 0 {
		return nil, false
	}

	current := n

	for _, key := range *keyChain {
		node, exists := current.findNodeByKey(key)

		if !exists {
			return nil, false
		}

		current = node
	}

	return current, true
}

func (n *node) findNodeByKey(key string) (*node, bool) {
	for _, child := range n.children {
		if child.key == key {
			return child, true
		}
	}

	return nil, false
}

func (n *node) appendNode(child *node) bool {
	if _, ok := n.findNodeByKey(child.key); ok {
		return false
	}

	n.children = append(n.children, child)

	return true
}
