package envh

import (
	"regexp"
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

	nodes := n.findAllChildsByKey("TEST2", false)

	assert.Len(t, *nodes, 2, "Must contains 2 elements")
}
