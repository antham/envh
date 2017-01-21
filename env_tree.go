package envh

import (
	"regexp"
	"strings"
)

func createTreeFromDelimiterFilteringByRegexp(reg *regexp.Regexp, delimiter string) (*node, error) {
	rootNode := newRootNode()

	for key, value := range *parseVars() {
		if reg.MatchString(key) {
			current := rootNode

			for _, component := range strings.Split(key, delimiter) {
				n, exists := current.findChildByKey(component)

				if exists {
					current = n
				} else {
					child := newNode()
					child.key = component
					current.appendChild(child)
					current = child
				}
			}

			if current.hasValue {
				return nil, ErrDuplicated
			}

			current.hasValue = true
			current.value = value
		}
	}

	return rootNode, nil
}

func getChildValueByKeyChain(node *node, keyChain *[]string) func() (string, bool) {
	return func() (string, bool) {
		n, exists := node.findChildByKeyChain(keyChain)

		if !exists {
			return "", false
		}

		return n.value, true
	}
}
