package envh

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateANode(t *testing.T) {
	n := newNode()

	assert.Equal(t, *n, node{childs: []*node{}}, "Must creates a new node")

	rootNode := newRootNode()

	assert.Equal(t, *rootNode, node{childs: []*node{}, root: true}, "Must creates a new root node")
}

func TestFindNodeByKey(t *testing.T) {
	root := newRootNode()

	node := newNode()
	node.key = "test"
	node.value = "value"
	root.appendNode(node)

	result, exists := root.findNodeByKey("test")

	assert.True(t, exists, "Must return true cause element was found")
	assert.Equal(t, node, result, "Must return child node with key test")

	_, exists = root.findNodeByKey("test1")

	assert.False(t, exists, "Must return false cause element was not found")
}

func TestAppendNode(t *testing.T) {
	root := newRootNode()

	node := newNode()
	node.key = "test"
	node.value = "value"

	result := root.appendNode(node)

	assert.True(t, result, "Must return true cause element was successfully added")
	assert.Equal(t, node, root.childs[0], "Must have node added as child")

	node2 := newNode()
	node2.key = "test"
	node2.value = "value2"

	result = root.appendNode(node2)

	assert.False(t, result, "Must return false cause an element with this key already exists")
	assert.Len(t, root.childs, 1, "Must still have one node")
	assert.Equal(t, node, root.childs[0], "Must have node added before")
}

func TestFindAllNodesByKey(t *testing.T) {
	nodes := map[string]*node{}

	root := newRootNode()
	n := root

	var accumulatedKey string

	for _, i := range []string{"1", "2", "3"} {
		t := newNode()
		t.key = "test" + i
		t.value = "value" + i
		n.appendNode(t)

		n = t

		if len(accumulatedKey) == 0 {
			accumulatedKey = i
		} else {
			accumulatedKey += "." + i
		}

		nodes[accumulatedKey] = t
	}

	accumulatedKey = ""
	n = root

	for _, i := range []string{"4", "5", "6", "3"} {
		t := newNode()
		t.key = "test" + i
		t.value = "value" + i
		n.appendNode(t)

		n = t

		if len(accumulatedKey) == 0 {
			accumulatedKey = i
		} else {
			accumulatedKey += "." + i
		}

		nodes[accumulatedKey] = t
	}

	results := root.findAllNodesByKey("test3", false)

	assert.Equal(t, []*node{nodes["4.5.6.3"], nodes["1.2.3"]}, *results, "Must recurse over tree to find keys")
}

func TestFindNodeByKeyChain(t *testing.T) {
	setTestingEnvsForTree()

	n := createTreeFromDelimiterFilteringByRegexp(regexp.MustCompile("ENVH"), "_")

	node, exists := n.findNodeByKeyChain(&[]string{"ENVH", "TEST1", "TEST5", "TEST6"})

	assert.True(t, exists, "Must find a node from this key chain")
	assert.Equal(t, "test3", node.value, "Must return correct node")

	for _, keyChain := range [][]string{
		[]string{},
		[]string{"ENV"},
		[]string{"ENVH", "TEST1", "TEST7", "TEST8"},
		[]string{"ENVH", "TEST1", "TEST6", "TEST8"},
	} {
		_, exists := n.findNodeByKeyChain(&keyChain)

		assert.False(t, exists, "Must not find a node from this key chain")
	}
}
