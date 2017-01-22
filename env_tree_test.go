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
		nodes := n.findAllChildsByKey(key, true)

		assert.Len(t, *nodes, 1, "Must contains 1 element")
		assert.Equal(t, expected, (*nodes)[0].value, "Must have correct value")
	}
}

func TestCreateTreeFromDelimiterFilteringByRegexpAndFindAllKeysWithAKey(t *testing.T) {
	setTestingEnvsForTree()

	n := createTreeFromDelimiterFilteringByRegexp(regexp.MustCompile("ENVH"), "_")

	nodes := n.findAllChildsByKey("TEST2", false)

	assert.Len(t, *nodes, 2, "Must contains 2 elements")
}

func TestNewEnvTreeWithAnInvalidRegexp(t *testing.T) {
	setTestingEnvsForTree()

	_, err := NewEnvTree("**", "_")

	assert.EqualError(t, err, "error parsing regexp: missing argument to repetition operator: `*`", "Must return an error when regexp is invalid")
}

func rebuildKeys(n *node, tmp []string, keys *[]string) {
	for _, child := range (*n).childs {
		t := append(tmp, child.key)

		rebuildKeys(child, t, keys)
	}

	if n.hasValue {
		*keys = append(*keys, strings.Join(tmp, "_"))
	}
}

func TestNewEnvTree(t *testing.T) {
	setTestingEnvsForTree()

	tree, _ := NewEnvTree("ENVH", "_")

	result := []string{}

	rebuildKeys(tree.root, []string{}, &result)

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
