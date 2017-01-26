package envh

import (
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTreeFromDelimiterFilteringByRegexp(t *testing.T) {
	setTestingEnvsForTree()

	n := createTreeFromDelimiterFilteringByRegexp(regexp.MustCompile("ENVH"), "_")

	for key, expected := range map[string]string{"TEST3": "test1", "TEST4": "test2", "TEST6": "test3", "TEST1": "test5", "TEST2": "test4"} {
		nodes := n.findAllNodesByKey(key, true)

		assert.Len(t, *nodes, 1, "Must contains 1 element")
		assert.Equal(t, expected, (*nodes)[0].value, "Must have correct value")
	}
}

func TestCreateTreeFromDelimiterFilteringByRegexpAndFindAllKeysWithAKey(t *testing.T) {
	setTestingEnvsForTree()

	n := createTreeFromDelimiterFilteringByRegexp(regexp.MustCompile("ENVH"), "_")

	nodes := n.findAllNodesByKey("TEST2", false)

	assert.Len(t, *nodes, 2, "Must contains 2 elements")
}

func TestNewEnvTreeWithAnInvalidRegexp(t *testing.T) {
	setTestingEnvsForTree()

	_, err := NewEnvTree("**", "_")

	assert.EqualError(t, err, "error parsing regexp: missing argument to repetition operator: `*`", "Must return an error when regexp is invalid")
}

func rebuildKeys(n *node, tmp []string, keys *[]string) {
	for _, child := range (*n).children {
		t := append(tmp, child.key)

		rebuildKeys(child, t, keys)
	}

	if n.hasValue {
		*keys = append(*keys, strings.Join(tmp, "_"))
	}
}

func TestNewEnvTree(t *testing.T) {
	setTestingEnvsForTree()

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	result := []string{}

	rebuildKeys(envTree.root, []string{}, &result)

	expected := []string{
		"ENVH_TEST1_TEST5_TEST6",
		"ENVH_TEST1_TEST2_TEST4",
		"ENVH_TEST1_TEST2_TEST3",
		"ENVH_TEST1_TEST7_TEST2",
		"ENVH_TEST1",
	}

	sort.Strings(expected)
	sort.Strings(result)

	assert.Equal(t, expected, result, "Must store all environment variables starting with envh in a tree")
}

func TestGetStringFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_TEST3", "test1")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	value, err := envTree.GetString("ENVH", "TEST1", "TEST2", "TEST3")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "test1", value, "Must return value")

	value, err = envTree.GetString("ENVH_TEST1000")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, "", value, "Must return empty string")
}

func TestGetIntFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_INT", "1")
	setEnv("ENVH_TEST1_TEST2_STRING", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	value, err := envTree.GetInt("ENVH", "TEST1", "TEST2", "INT")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, 1, value, "Must return value")

	value, err = envTree.GetInt("TEST100")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, 0, value, "Must return value")

	value, err = envTree.GetInt("ENVH", "TEST1", "TEST2", "STRING")

	assert.EqualError(t, err, "Variable can't be converted", "Must return an error when variable can't be converted")
	assert.Equal(t, 0, value, "Must return empty string")
}

func TestGetBoolFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_BOOL", "1")
	setEnv("ENVH_TEST1_TEST2_STRING", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	value, err := envTree.GetBool("ENVH", "TEST1", "TEST2", "BOOL")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, true, value, "Must return value")

	value, err = envTree.GetBool("TEST100")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, false, value, "Must return value")

	value, err = envTree.GetBool("ENVH", "TEST1", "TEST2", "STRING")

	assert.EqualError(t, err, "Variable can't be converted", "Must return an error when variable can't be converted")
	assert.Equal(t, false, value, "Must return empty string")
}

func TestGetFloatFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_FLOAT", "0.01")
	setEnv("ENVH_TEST1_TEST2_STRING", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	value, err := envTree.GetFloat("ENVH", "TEST1", "TEST2", "FLOAT")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, float32(0.01), value, "Must return value")

	value, err = envTree.GetFloat("TEST100")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, float32(0), value, "Must return value")

	value, err = envTree.GetFloat("ENVH", "TEST1", "TEST2", "STRING")

	assert.EqualError(t, err, "Variable can't be converted", "Must return an error when variable can't be converted")
	assert.Equal(t, float32(0), value, "Must return empty string")
}

func TestExistsFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_TEST3", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	assert.True(t, envTree.Exists("ENVH", "TEST1", "TEST2", "TEST3"), "Must return true if environment node exists")

	assert.False(t, envTree.Exists("ENVH", "TEST1", "TEST2", "TEST10000"), "Must return false if environment node doesn't exist")
}

func TestHasValueFromTree(t *testing.T) {
	setEnv("ENVH_TEST10_TEST20_TEST30", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	v, err := envTree.HasValue("ENVH", "TEST10", "TEST20", "TEST30")

	assert.NoError(t, err, "Must returns no error")
	assert.True(t, v, "Must return true if environment node has a value")

	v, err = envTree.HasValue("ENVH", "TEST10", "TEST20")

	assert.NoError(t, err, "Must returns no error")
	assert.False(t, v, "Must return false if environment node doesn't have value")

	_, err = envTree.HasValue("ENVH", "TEST10", "TEST20", "TEST10000")

	assert.EqualError(t, err, ErrNodeNotFound.Error(), "Must returns an error, node doesn't exists")
}

func TestGetChildrenKeysFromTree(t *testing.T) {
	setEnv("ENVH_TEST11_TEST12_TEST13_TEST14", "test1")
	setEnv("ENVH_TEST11_TEST12_TEST13_TEST15", "test2")
	setEnv("ENVH_TEST11_TEST12_TEST13_TEST16", "test3")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	v, err := envTree.GetChildrenKeys("ENVH", "TEST11", "TEST12", "TEST13")

	sort.Strings(v)

	assert.NoError(t, err, "Must returns no error")
	assert.Equal(t, []string{"TEST14", "TEST15", "TEST16"}, v, "Must return all node children keys")

	v, err = envTree.GetChildrenKeys("ENVH", "TEST11", "TEST12", "TEST13", "TEST14")

	assert.NoError(t, err, "Must returns no error")
	assert.Empty(t, v, "Must return no node children keys")

	_, err = envTree.GetChildrenKeys("ENVH", "TEST11", "TEST12", "TEST13", "TEST10000")

	assert.EqualError(t, err, ErrNodeNotFound.Error(), "Must returns an error, node doesn't exists")
}
