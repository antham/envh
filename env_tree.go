package envh

import (
	"regexp"
	"strings"
)

func createTreeFromDelimiterFilteringByRegexp(reg *regexp.Regexp, delimiter string) *node {
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

			current.hasValue = true
			current.value = value
		}
	}

	return rootNode
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
