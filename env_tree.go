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

// FindString returns a string if key chain exists
// or an error otherwise
func (e EnvTree) FindString(keyChain ...string) (string, error) {
	return getString(getNodeValueByKeyChain(e.root, &keyChain))
}

// FindInt returns an integer if key chain exists
// or an error if value is not an integer or doesn't exist
func (e EnvTree) FindInt(keyChain ...string) (int, error) {
	return getInt(getNodeValueByKeyChain(e.root, &keyChain))
}

// FindFloat returns a float if key chain exists
// or an error if value is not a float or doesn't exist
func (e EnvTree) FindFloat(keyChain ...string) (float32, error) {
	return getFloat(getNodeValueByKeyChain(e.root, &keyChain))
}

// FindBool returns a boolean if key chain exists
// or an error if value is not a boolean or doesn't exist
func (e EnvTree) FindBool(keyChain ...string) (bool, error) {
	return getBool(getNodeValueByKeyChain(e.root, &keyChain))
}

// IsExistingSubTree returns true if key chain has a tree associated or false if not
func (e EnvTree) IsExistingSubTree(keyChain ...string) bool {
	_, exists := e.root.findNodeByKeyChain(&keyChain)

	return exists
}

// HasSubTreeValue returns true if key chain has a value or false if not.
// If sub node doesn't exist, it returns an error ErrNodeNotFound
// as second value
func (e EnvTree) HasSubTreeValue(keyChain ...string) (bool, error) {
	n, exists := e.root.findNodeByKeyChain(&keyChain)

	if !exists {
		return false, NodeNotFoundError{keyChain}
	}

	return n.hasValue, nil
}

// FindSubTree returns underlying tree from key chain,
// for instance given A -> B -> C -> D tree,
// "A" "B" "C" key chain will return C sub tree.
// If no node is found, it returns an error ErrNodeNotFound as
// second value
func (e EnvTree) FindSubTree(keyChain ...string) (EnvTree, error) {
	if n, exists := e.root.findNodeByKeyChain(&keyChain); exists {
		return EnvTree{n}, nil
	}

	return EnvTree{}, NodeNotFoundError{keyChain}
}

// FindChildrenKeys returns all children keys for a given key chain.
// If sub node doesn't exist, it returns an error ErrNodeNotFound
// as second value
func (e EnvTree) FindChildrenKeys(keyChain ...string) ([]string, error) {
	n, exists := e.root.findNodeByKeyChain(&keyChain)

	if !exists {
		return []string{}, NodeNotFoundError{keyChain}
	}

	keys := []string{}

	for _, c := range n.children {
		keys = append(keys, c.key)
	}

	return keys, nil
}

// GetChildrenKeys retrieves all current tree children node keys
func (e EnvTree) GetChildrenKeys() []string {
	keys := []string{}

	for _, c := range e.root.children {
		keys = append(keys, c.key)
	}

	return keys
}

// GetString returns current tree value as string if value exists
// or an error as second parameter
func (e EnvTree) GetString() (string, error) {
	return getString(e.getValue())
}

// GetInt returns current tree value as int if value exists
// or an error if value is not an integer or doesn't exist
func (e EnvTree) GetInt() (int, error) {
	return getInt(e.getValue())
}

// GetFloat returns current tree value as float if value exists
// or an error if value is not a float or doesn't exist
func (e EnvTree) GetFloat() (float32, error) {
	return getFloat(e.getValue())
}

// GetBool returns current tree value as boolean if value exists
// or an error if value is not a boolean or doesn't exist
func (e EnvTree) GetBool() (bool, error) {
	return getBool(e.getValue())
}

// HasValue returns true if current tree has a value defined
// false otherwise
func (e EnvTree) HasValue() bool {
	return e.root.hasValue
}

// GetKey returns current tree key
func (e EnvTree) GetKey() string {
	return e.root.key
}

// PopulateStruct fills a structure with datas extracted.
// Missing values are ignored and only type errors are reported.
// It's possible to control the way struct fields are defined
// implementing StructWalker interface on structure,
// checkout StructWalker documentation for further examples.
func (e EnvTree) PopulateStruct(structure interface{}) error {
	return populateStructFromEnvTree(structure, &e, false)
}

// PopulateStructWithStrictMode fills a structure with datas extracted.
// A missing environment variable returns an error and type errors are reported.
// It's possible to control the way struct fields are defined
// implementing StructWalker interface on structure,
// checkout StructWalker documentation for further examples.
func (e EnvTree) PopulateStructWithStrictMode(structure interface{}) error {
	return populateStructFromEnvTree(structure, &e, true)
}

func (e EnvTree) getValue() func() (string, bool) {
	return func() (string, bool) {
		if e.root.hasValue {
			return e.root.value, true
		}

		return "", false
	}
}

func createTreeFromDelimiterFilteringByRegexp(reg *regexp.Regexp, delimiter string) *node {
	rootNode := newNode()

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
