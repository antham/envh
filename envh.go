package envh

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ErrNotFound is triggered when environment variable cannot be found
var ErrNotFound = fmt.Errorf("Variable not found")

// ErrWrongType is triggered when we try to convert variable to a wrong type
var ErrWrongType = fmt.Errorf("Variable can't be converted")

// ErrDuplicated is triggered when a variable is already defined in a tree structure
var ErrDuplicated = fmt.Errorf("Variable was already defined before")

func parseVars() *map[string]string {
	results := map[string]string{}

	for _, v := range os.Environ() {
		e := strings.SplitN(v, "=", 2)

		results[e[0]] = e[1]
	}

	return &results
}

func getString(fun func() (string, bool)) (string, error) {
	if v, ok := fun(); ok {
		return v, nil
	}

	return "", ErrNotFound
}

func getInt(fun func() (string, bool)) (int, error) {
	v, ok := fun()

	if !ok {
		return 0, ErrNotFound
	}

	i, err := strconv.Atoi(v)

	if err != nil {
		return 0, ErrWrongType
	}

	return i, nil
}

func getFloat(fun func() (string, bool)) (float32, error) {
	v, ok := fun()

	if !ok {
		return 0, ErrNotFound
	}

	f, err := strconv.ParseFloat(v, 32)

	if err != nil {
		return 0, ErrWrongType
	}

	return float32(f), nil
}

func getBool(fun func() (string, bool)) (bool, error) {
	v, ok := fun()

	if !ok {
		return false, ErrNotFound
	}

	b, err := strconv.ParseBool(v)

	if err != nil {
		return false, ErrWrongType
	}

	return b, nil
}

// Env manage environment variables
// by giving high level api to interact with them
type Env struct {
	envs *map[string]string
}

// NewEnv creates a new Env instance
func NewEnv() Env {
	return Env{parseVars()}
}

// GetAllValues retrieves a slice of all environment variables values
func (e Env) GetAllValues() []string {
	results := []string{}

	for _, v := range *e.envs {
		results = append(results, v)
	}

	return results
}

// GetAllKeys retrieves a slice of all environment variables keys
func (e Env) GetAllKeys() []string {
	results := []string{}

	for k := range *e.envs {
		results = append(results, k)
	}

	return results
}

// GetString return a string if variable exists
// or an error otherwise
func (e Env) GetString(key string) (string, error) {
	return getString(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
}

// GetInt return an integer if variable exists
// or an error if value is not an integer or doesn't exist
func (e Env) GetInt(key string) (int, error) {
	return getInt(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
}

// GetFloat return a float if variable exists
// or an error if value is not a float or doesn't exist
func (e Env) GetFloat(key string) (float32, error) {
	return getFloat(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
}

// GetBool return a boolean if variable exists
// or an error if value is not a boolean or doesn't exist
func (e Env) GetBool(key string) (bool, error) {
	return getBool(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
}

// FindEntries retrieves all keys matching a given regexp and their
// corresponding values
func (e Env) FindEntries(reg string) (map[string]string, error) {
	results := map[string]string{}

	r, err := regexp.Compile(reg)

	if err != nil {
		return results, err
	}

	for k, v := range *e.envs {
		if r.MatchString(k) {
			results[k] = v
		}
	}

	return results, nil
}

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



}
