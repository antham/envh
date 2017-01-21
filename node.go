package envh

type node struct {
	childs   []*node
	key      string
	value    string
	hasValue bool
	root     bool
}

func newNode() *node {
	return &node{childs: []*node{}}
}

func newRootNode() *node {
	return &node{childs: []*node{}, root: true}
}

func (n *node) findAllChildsByKey(key string, withValue bool) *[]*node {
	results := []*node{}
	nodes := n.childs

	for {
		carry := []*node{}

		for _, node := range nodes {
			if node.key == key {
				if withValue && node.hasValue || !withValue {
					results = append(results, node)
				}
			}

			carry = append(carry, node.childs...)
		}

		nodes = carry

		if len(carry) == 0 {
			return &results
		}
	}
}

func (n *node) findChildByKeyChain(keyChain *[]string) (*node, bool) {
	if len(*keyChain) == 0 {
		return nil, false
	}

	current := n

	for _, key := range *keyChain {
		node, exists := current.findChildByKey(key)

		if !exists {
			return nil, false
		}

		current = node
	}

	return current, true
}

func (n *node) findChildByKey(key string) (*node, bool) {
	for _, child := range n.childs {
		if child.key == key {
			return child, true
		}
	}

	return nil, false
}

func (n *node) appendChild(child *node) bool {
	if _, ok := n.findChildByKey(child.key); ok {
		return false
	}

	n.childs = append(n.childs, child)

	return true
}
