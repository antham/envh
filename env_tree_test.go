package envh

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTreeFromDelimiterFilteringByRegexp(t *testing.T) {
	setTestingEnvsForTree()

	n, err := createTreeFromDelimiterFilteringByRegexp(regexp.MustCompile("ENVH"), "_")

	nodes := n.findAllChildsByKey("TEST3", true)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *nodes, 1, "Must contains 1 element")
	assert.Equal(t, "test1", (*nodes)[0].value, "Must have correct value")

	nodes = n.findAllChildsByKey("TEST4", true)

	assert.Len(t, *nodes, 1, "Must contains 1 element")
	assert.Equal(t, "test2", (*nodes)[0].value, "Must have correct value")

	nodes = n.findAllChildsByKey("TEST6", true)

	assert.Len(t, *nodes, 1, "Must contains 1 element")
	assert.Equal(t, "test3", (*nodes)[0].value, "Must have correct value")

	nodes = n.findAllChildsByKey("TEST1", true)

	assert.Len(t, *nodes, 1, "Must contains 1 element")
	assert.Equal(t, "test5", (*nodes)[0].value, "Must have correct value")

	// Find all childs with a value defined
	nodes = n.findAllChildsByKey("TEST2", true)

	assert.Len(t, *nodes, 1, "Must contains 1 element")

	assert.Equal(t, "test4", (*nodes)[0].value, "Must have correct value")

	// Find all childs with values or not
	nodes = n.findAllChildsByKey("TEST2", false)

	assert.Len(t, *nodes, 2, "Must contains 2 elements")
}
