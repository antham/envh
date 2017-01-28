package envh

import (
	"regexp"
	"strings"
)

// EnvTree manages environment variables through a tree structure
// to store a config the same way as in a yaml file or whatever
// format allows to store a config hierarchically
type EnvTree struct {
	root *node
}

// NewEnvTree creates an environment variable tree.
// A delimiter is used to split key, reg is a regexp
// used to filter entries
func NewEnvTree(reg string, delimiter string) (EnvTree, error) {
	r, err := regexp.Compile(reg)

	if err != nil {
		return EnvTree{}, err
	}

	t := createTreeFromDelimiterFilteringByRegexp(r, delimiter)

	return EnvTree{t}, nil
}

// GetString returns a string if variable exists
// or an error otherwise
func (e EnvTree) GetString(keyChain ...string) (string, error) {
	return getString(getNodeValueByKeyChain(e.root, &keyChain))
}

// GetInt returns an integer if variable exists
// or an error if value is not an integer or doesn't exist
func (e EnvTree) GetInt(keyChain ...string) (int, error) {
	return getInt(getNodeValueByKeyChain(e.root, &keyChain))
}

// GetFloat returns a float if variable exists
// or an error if value is not a float or doesn't exist
func (e EnvTree) GetFloat(keyChain ...string) (float32, error) {
	return getFloat(getNodeValueByKeyChain(e.root, &keyChain))
}

// GetBool returns a boolean if variable exists
// or an error if value is not a boolean or doesn't exist
func (e EnvTree) GetBool(keyChain ...string) (bool, error) {
	return getBool(getNodeValueByKeyChain(e.root, &keyChain))
}

// Exists returns true if sub node exits or false if not
func (e EnvTree) Exists(keyChain ...string) bool {
	_, exists := e.root.findNodeByKeyChain(&keyChain)

	return exists
}

// HasValue returns true if sub node has a value or false if not.
// If sub node doesn't exist, it returns an error ErrNodeNotFound
// as second value
func (e EnvTree) HasValue(keyChain ...string) (bool, error) {
	n, exists := e.root.findNodeByKeyChain(&keyChain)

	if !exists {
		return false, ErrNodeNotFound
	}

	return n.hasValue, nil
}

// GetSubTree returns underlying tree from key chain,
// for instance given A -> B -> C -> D tree,
// "A" "B" "C" key chain will return C sub tree.
// If no node is found, it returns an error ErrNodeNotFound as
// second value
func (e EnvTree) GetSubTree(keyChain ...string) (EnvTree, error) {
	if n, exists := e.root.findNodeByKeyChain(&keyChain); exists {
		return EnvTree{n}, nil
	}

	return EnvTree{}, ErrNodeNotFound
}

// GetChildrenKeys returns all children keys for a given key chain.
// If sub node doesn't exist, it returns an error ErrNodeNotFound
// as second value
func (e EnvTree) GetChildrenKeys(keyChain ...string) ([]string, error) {
	n, exists := e.root.findNodeByKeyChain(&keyChain)

	if !exists {
		return []string{}, ErrNodeNotFound
	}

	keys := []string{}

	for _, c := range n.children {
		keys = append(keys, c.key)
	}

	return keys, nil
}

func createTreeFromDelimiterFilteringByRegexp(reg *regexp.Regexp, delimiter string) *node {
	rootNode := newRootNode()

	for key, value := range *parseVars() {
		if reg.MatchString(key) {
			current := rootNode

			for _, component := range strings.Split(key, delimiter) {
				n, exists := current.findNodeByKey(component)

				if exists {
					current = n
				} else {
					child := newNode()
					child.key = component
					current.appendNode(child)
					current = child
				}
			}

			current.hasValue = true
			current.value = value
		}
	}

	return rootNode
}

func getNodeValueByKeyChain(node *node, keyChain *[]string) func() (string, bool) {
	return func() (string, bool) {
		n, exists := node.findNodeByKeyChain(keyChain)

		if !exists {
			return "", false
		}

		return n.value, true
	}
}
