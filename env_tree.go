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

// FindStringUnsecured is insecured version of FindString to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing, it returns default zero string value.
// This function has to be used carefully
func (e EnvTree) FindStringUnsecured(keyChain ...string) string {
	if val, err := getString(getNodeValueByKeyChain(e.root, &keyChain)); err == nil {
		return val
	}

	return ""
}

// FindInt returns an integer if key chain exists
// or an error if value is not an integer or doesn't exist
func (e EnvTree) FindInt(keyChain ...string) (int, error) {
	return getInt(getNodeValueByKeyChain(e.root, &keyChain))
}

// FindIntUnsecured is insecured version of FindInt to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing or not an int value, it returns default zero int value.
// This function has to be used carefully
func (e EnvTree) FindIntUnsecured(keyChain ...string) int {
	if val, err := getInt(getNodeValueByKeyChain(e.root, &keyChain)); err == nil {
		return val
	}

	return 0
}

// FindFloat returns a float if key chain exists
// or an error if value is not a float or doesn't exist
func (e EnvTree) FindFloat(keyChain ...string) (float32, error) {
	return getFloat(getNodeValueByKeyChain(e.root, &keyChain))
}

// FindFloatUnsecured is insecured version of FindFloat to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing or not a floating value, it returns default zero floating value.
// This function has to be used carefully
func (e EnvTree) FindFloatUnsecured(keyChain ...string) float32 {
	if val, err := getFloat(getNodeValueByKeyChain(e.root, &keyChain)); err == nil {
		return val
	}

	return 0
}

// FindBool returns a boolean if key chain exists
// or an error if value is not a boolean or doesn't exist
func (e EnvTree) FindBool(keyChain ...string) (bool, error) {
	return getBool(getNodeValueByKeyChain(e.root, &keyChain))
}

// FindBoolUnsecured is insecured version of FindBool to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing or not a boolean value, it returns default zero boolean value.
// This function has to be used carefully
func (e EnvTree) FindBoolUnsecured(keyChain ...string) bool {
	if val, err := getBool(getNodeValueByKeyChain(e.root, &keyChain)); err == nil {
		return val
	}

	return false
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

// HasSubTreeValueUnsecured is insecured version of HasSubTreeValue to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the node doesn't exist, it returns false.
// This function has to be used carefully
func (e EnvTree) HasSubTreeValueUnsecured(keyChain ...string) bool {
	n, exists := e.root.findNodeByKeyChain(&keyChain)

	if !exists {
		return false
	}

	return n.hasValue
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

// FindSubTreeUnsecured is insecured version of FindSubTree to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the node doesn't exist, it returns empty EnvTree.
// This function has to be used carefully
func (e EnvTree) FindSubTreeUnsecured(keyChain ...string) EnvTree {
	if n, exists := e.root.findNodeByKeyChain(&keyChain); exists {
		return EnvTree{n}
	}

	return EnvTree{}
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

// FindChildrenKeysUnsecured is insecured version of FindChildrenKeys to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the node doesn't exist, it returns empty string slice.
// This function has to be used carefully
func (e EnvTree) FindChildrenKeysUnsecured(keyChain ...string) []string {
	n, exists := e.root.findNodeByKeyChain(&keyChain)

	if !exists {
		return []string{}
	}

	keys := []string{}

	for _, c := range n.children {
		keys = append(keys, c.key)
	}

	return keys
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
	return getString(getRootValue(e))
}

// GetStringUnsecured is insecured version of GetString to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing, it returns default zero string value.
// This function has to be used carefully
func (e EnvTree) GetStringUnsecured() string {
	if val, err := getString(getRootValue(e)); err == nil {
		return val
	}

	return ""
}

// GetInt returns current tree value as int if value exists
// or an error if value is not an integer or doesn't exist
func (e EnvTree) GetInt() (int, error) {
	return getInt(getRootValue(e))
}

// GetIntUnsecured is insecured version of GetInt to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing or not an int value, it returns default zero int value.
// This function has to be used carefully
func (e EnvTree) GetIntUnsecured() int {
	if val, err := getInt(getRootValue(e)); err == nil {
		return val
	}

	return 0
}

// GetFloat returns current tree value as float if value exists
// or an error if value is not a float or doesn't exist
func (e EnvTree) GetFloat() (float32, error) {
	return getFloat(getRootValue(e))
}

// GetFloatUnsecured is insecured version of GetFloat to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing or not a floating value, it returns default zero floating value.
// This function has to be used carefully
func (e EnvTree) GetFloatUnsecured() float32 {
	if val, err := getFloat(getRootValue(e)); err == nil {
		return val
	}

	return 0
}

// GetBool returns current tree value as boolean if value exists
// or an error if value is not a boolean or doesn't exist
func (e EnvTree) GetBool() (bool, error) {
	return getBool(getRootValue(e))
}

// GetBoolUnsecured is insecured version of GetBool to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing or not a boolean value, it returns default zero boolean value.
// This function has to be used carefully
func (e EnvTree) GetBoolUnsecured() bool {
	if val, err := getBool(getRootValue(e)); err == nil {
		return val
	}

	return false
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

func getRootValue(tree EnvTree) func() (string, bool) {
	return func() (string, bool) {
		if tree.root.hasValue {
			return tree.root.value, true
		}

		return "", false
	}
}

func createBranch(key string, value string, delimiter string, current *node) {
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

func createTreeFromDelimiterFilteringByRegexp(reg *regexp.Regexp, delimiter string) *node {
	rootNode := newNode()

	for key, value := range *parseVars() {
		if reg.MatchString(key) {
			createBranch(key, value, delimiter, rootNode)
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
