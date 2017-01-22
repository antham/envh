package envh

import (
	"regexp"
	"strings"
)

// EnvTree manage environment variables through a tree structure
// to store a config the same way as in a yaml file or whatever
// format allows to store a config hierarchically
type EnvTree struct {
	root *node
}

// NewEnvTree creates an environment tree
// delimiter is used to split key, reg is a regexp used to filter entries
func NewEnvTree(reg string, delimiter string) (EnvTree, error) {
	r, err := regexp.Compile(reg)

	if err != nil {
		return EnvTree{}, err
	}

	t := createTreeFromDelimiterFilteringByRegexp(r, delimiter)

	return EnvTree{t}, nil
}

// GetString return a string if variable exists
// or an error otherwise
func (e EnvTree) GetString(keyChain ...string) (string, error) {
	return getString(getChildValueByKeyChain(e.root, &keyChain))
}

// GetInt return an integer if variable exists
// or an error if value is not an integer or doesn't exist
func (e EnvTree) GetInt(keyChain ...string) (int, error) {
	return getInt(getChildValueByKeyChain(e.root, &keyChain))
}

// GetFloat return a float if variable exists
// or an error if value is not a float or doesn't exist
func (e EnvTree) GetFloat(keyChain ...string) (float32, error) {
	return getFloat(getChildValueByKeyChain(e.root, &keyChain))
}

// GetBool return a boolean if variable exists
// or an error if value is not a boolean or doesn't exist
func (e EnvTree) GetBool(keyChain ...string) (bool, error) {
	return getBool(getChildValueByKeyChain(e.root, &keyChain))
}

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
