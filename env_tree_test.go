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

func TestFindStringFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_TEST3", "test1")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	value, err := envTree.FindString("ENVH", "TEST1", "TEST2", "TEST3")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "test1", value, "Must return value")

	value, err = envTree.FindString("ENVH_TEST1000")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, "", value, "Must return empty string")
}

func TestFindIntFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_INT", "1")
	setEnv("ENVH_TEST1_TEST2_STRING", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	value, err := envTree.FindInt("ENVH", "TEST1", "TEST2", "INT")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, 1, value, "Must return value")

	value, err = envTree.FindInt("TEST100")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, 0, value, "Must return value")

	value, err = envTree.FindInt("ENVH", "TEST1", "TEST2", "STRING")

	assert.EqualError(t, err, `Value "test" can't be converted to type "int"`, "Must return an error when variable can't be converted")
	assert.Equal(t, 0, value, "Must return empty string")
}

func TestFindBoolFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_BOOL", "1")
	setEnv("ENVH_TEST1_TEST2_STRING", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	value, err := envTree.FindBool("ENVH", "TEST1", "TEST2", "BOOL")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, true, value, "Must return value")

	value, err = envTree.FindBool("TEST100")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, false, value, "Must return value")

	value, err = envTree.FindBool("ENVH", "TEST1", "TEST2", "STRING")

	assert.EqualError(t, err, `Value "test" can't be converted to type "bool"`, "Must return an error when variable can't be converted")
	assert.Equal(t, false, value, "Must return empty string")
}

func TestFindFloatFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_FLOAT", "0.01")
	setEnv("ENVH_TEST1_TEST2_STRING", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	value, err := envTree.FindFloat("ENVH", "TEST1", "TEST2", "FLOAT")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, float32(0.01), value, "Must return value")

	value, err = envTree.FindFloat("TEST100")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, float32(0), value, "Must return value")

	value, err = envTree.FindFloat("ENVH", "TEST1", "TEST2", "STRING")

	assert.EqualError(t, err, `Value "test" can't be converted to type "float"`, "Must return an error when variable can't be converted")
	assert.Equal(t, float32(0), value, "Must return empty string")
}

func TestIsExistingSubTreeFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_TEST3", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	assert.True(t, envTree.IsExistingSubTree("ENVH", "TEST1", "TEST2", "TEST3"), "Must return true if environment node exists")

	assert.False(t, envTree.IsExistingSubTree("ENVH", "TEST1", "TEST2", "TEST10000"), "Must return false if environment node doesn't exist")
}

func TestHasSubTreeValueFromTree(t *testing.T) {
	setEnv("ENVH_TEST10_TEST20_TEST30", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	v, err := envTree.HasSubTreeValue("ENVH", "TEST10", "TEST20", "TEST30")

	assert.NoError(t, err, "Must returns no error")
	assert.True(t, v, "Must return true if environment node has a value")

	v, err = envTree.HasSubTreeValue("ENVH", "TEST10", "TEST20")

	assert.NoError(t, err, "Must returns no error")
	assert.False(t, v, "Must return false if environment node doesn't have value")

	_, err = envTree.HasSubTreeValue("ENVH", "TEST10", "TEST20", "TEST10000")

	assert.EqualError(t, err, `No node found at path "ENVH -> TEST10 -> TEST20 -> TEST10000"`, "Must returns an error, node doesn't exists")
}

func TestFindChildrenKeysFromTree(t *testing.T) {
	setEnv("ENVH_TEST11_TEST12_TEST13_TEST14", "test1")
	setEnv("ENVH_TEST11_TEST12_TEST13_TEST15", "test2")
	setEnv("ENVH_TEST11_TEST12_TEST13_TEST16", "test3")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	v, err := envTree.FindChildrenKeys("ENVH", "TEST11", "TEST12", "TEST13")

	sort.Strings(v)

	assert.NoError(t, err, "Must returns no error")
	assert.Equal(t, []string{"TEST14", "TEST15", "TEST16"}, v, "Must return all node children keys")

	v, err = envTree.FindChildrenKeys("ENVH", "TEST11", "TEST12", "TEST13", "TEST14")

	assert.NoError(t, err, "Must returns no error")
	assert.Empty(t, v, "Must return no node children keys")

	_, err = envTree.FindChildrenKeys("ENVH", "TEST11", "TEST12", "TEST13", "TEST10000")

	assert.EqualError(t, err, `No node found at path "ENVH -> TEST11 -> TEST12 -> TEST13 -> TEST10000"`, "Must returns an error, node doesn't exists")
}

func TestFindSubTreeFromTree(t *testing.T) {
	setEnv("ENVH_TEST11_TEST12_TEST13_TEST14", "test1")
	setEnv("ENVH_TEST11_TEST12_TEST13_TEST15", "test2")
	setEnv("ENVH_TEST11_TEST12_TEST13_TEST16", "test3")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	tree, err := envTree.FindSubTree("ENVH", "TEST11", "TEST12", "TEST13")

	assert.NoError(t, err, "Must returns no error")

	v, err := tree.FindString("TEST14")

	assert.NoError(t, err, "Must returns no error")

	assert.Equal(t, "test1", v, "Must return value corresponding to key TEST14")

	_, err = envTree.FindSubTree("ENVH", "TEST11", "TEST12", "TEST13", "TEST10000")

	assert.EqualError(t, err, `No node found at path "ENVH -> TEST11 -> TEST12 -> TEST13 -> TEST10000"`, "Must returns an error, node doesn't exists")
}

func TestGetStringFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_TEST3", "test1")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	subTree, err := envTree.FindSubTree("ENVH", "TEST1", "TEST2", "TEST3")

	assert.NoError(t, err, "Must returns no error")

	value, err := subTree.GetString()

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "test1", value, "Must return value")

	subTree, err = envTree.FindSubTree("ENVH", "TEST1", "TEST2")

	assert.NoError(t, err, "Must return no errors")

	value, err = subTree.GetString()

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, "", value, "Must return empty string")
}

func TestGetIntFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_INT", "1")
	setEnv("ENVH_TEST1_TEST2_STRING", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	subTree, err := envTree.FindSubTree("ENVH", "TEST1", "TEST2", "INT")

	assert.NoError(t, err, "Must return no errors")

	value, err := subTree.GetInt()

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, 1, value, "Must return value")

	subTree, err = envTree.FindSubTree("ENVH", "TEST1", "TEST2")

	assert.NoError(t, err, "Must return no errors")

	value, err = subTree.GetInt()

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, 0, value, "Must return value")

	subTree, err = envTree.FindSubTree("ENVH", "TEST1", "TEST2", "STRING")

	assert.NoError(t, err, "Must return no errors")

	value, err = subTree.GetInt()

	assert.EqualError(t, err, `Value "test" can't be converted to type "int"`, "Must return an error when variable can't be converted")
	assert.Equal(t, 0, value, "Must return empty string")
}

func TestGetBoolFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_BOOL", "true")
	setEnv("ENVH_TEST1_TEST2_STRING", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	subTree, err := envTree.FindSubTree("ENVH", "TEST1", "TEST2", "BOOL")

	assert.NoError(t, err, "Must return no errors")

	value, err := subTree.GetBool()

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, true, value, "Must return value")

	subTree, err = envTree.FindSubTree("ENVH", "TEST1", "TEST2")

	assert.NoError(t, err, "Must return no errors")

	value, err = subTree.GetBool()

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, false, value, "Must return value")

	subTree, err = envTree.FindSubTree("ENVH", "TEST1", "TEST2", "STRING")

	assert.NoError(t, err, "Must return no errors")

	value, err = subTree.GetBool()

	assert.EqualError(t, err, `Value "test" can't be converted to type "bool"`, "Must return an error when variable can't be converted")
	assert.Equal(t, false, value, "Must return empty string")
}

func TestGetFloatFromTree(t *testing.T) {
	setEnv("ENVH_TEST1_TEST2_FLOAT", "0.01")
	setEnv("ENVH_TEST1_TEST2_STRING", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no error")

	subTree, err := envTree.FindSubTree("ENVH", "TEST1", "TEST2", "FLOAT")

	assert.NoError(t, err, "Must return no errors")

	value, err := subTree.GetFloat()

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, float32(0.01), value, "Must return value")

	subTree, err = envTree.FindSubTree("ENVH", "TEST1", "TEST2")

	assert.NoError(t, err, "Must return no errors")

	value, err = subTree.GetFloat()

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, float32(0), value, "Must return value")

	subTree, err = envTree.FindSubTree("ENVH", "TEST1", "TEST2", "STRING")

	assert.NoError(t, err, "Must return no errors")

	value, err = subTree.GetFloat()

	assert.EqualError(t, err, `Value "test" can't be converted to type "float"`, "Must return an error when variable can't be converted")
	assert.Equal(t, float32(0), value, "Must return empty string")
}

func TestGetChildrenKeys(t *testing.T) {
	setEnv("KEY1", "test")
	setEnv("KEY2", "test")
	setEnv("KEY3", "test")

	envTree, err := NewEnvTree("^KEY[0-9]", " ")

	assert.NoError(t, err, "Must returns no errors")

	result := envTree.GetChildrenKeys()

	sort.Strings(result)

	assert.Equal(t, []string{"KEY1", "KEY2", "KEY3"}, result, "Must returns an array of all child keys")

	subTree, err := envTree.FindSubTree("KEY1")

	assert.NoError(t, err, "Must returns no errors")

	assert.Empty(t, subTree.GetChildrenKeys(), "Must returns an empty array of keys")
}

func TestHasValue(t *testing.T) {
	setEnv("ENVH_KEY1_KEY2", "test")

	envTree, err := NewEnvTree("ENVH", "_")

	assert.NoError(t, err, "Must returns no errors")

	subTree, err := envTree.FindSubTree("ENVH", "KEY1", "KEY2")

	assert.NoError(t, err, "Must returns no errors")

	result := subTree.HasValue()

	assert.True(t, result, "Must returns true, sub tree has value")

	subTree, err = envTree.FindSubTree("ENVH", "KEY1")

	assert.NoError(t, err, "Must returns no errors")

	result = subTree.HasValue()

	assert.False(t, result, "Must returns false, sub tree has value")
}

func TestGetKey(t *testing.T) {
	setEnv("KEY1", "test")

	envTree, err := NewEnvTree("KEY[0-9]", " ")

	assert.NoError(t, err, "Must returns no errors")

	result := envTree.GetKey()

	assert.Equal(t, result, "", "Must returns an empty string, original root node has no value")

	subTree, err := envTree.FindSubTree("KEY1")

	assert.NoError(t, err, "Must returns no errors")

	result = subTree.GetKey()

	assert.Equal(t, result, "KEY1", "Must returns key")
}

func TestPopulateStruct(t *testing.T) {
	setEnv("TEST_WHATEVER", "string")

	type TEST struct {
		WHATEVER string
	}

	envTree, err := NewEnvTree("TEST", "_")

	assert.NoError(t, err)

	actual := TEST{}

	err = envTree.PopulateStruct(&actual)

	assert.NoError(t, err)
	assert.Equal(t, TEST{"string"}, actual)

	restoreEnvs()
}

func TestPopulateStructWithStrictModeWithAnError(t *testing.T) {
	setEnv("TEST_STRING", "string")

	type TEST struct {
		WHATEVER string
	}

	envTree, err := NewEnvTree("TEST", "_")

	assert.NoError(t, err)

	actual := TEST{}

	err = envTree.PopulateStructWithStrictMode(&actual)

	assert.EqualError(t, err, "Variable not found")

	restoreEnvs()
}
