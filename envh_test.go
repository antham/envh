package envh

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setTestingEnvs() {
	datas := map[string]string{
		"TEST1": "test1",
		"TEST2": "=test2=",
	}

	for k, v := range datas {
		err := os.Setenv(k, v)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func TestParseVars(t *testing.T) {
	setTestingEnvs()
	result := parseVars()

	assert.Equal(t, "test1", (*result)["TEST1"], "Must extract and parse environment variables")
	assert.Contains(t, "=test2=", (*result)["TEST2"], "Must extract and parse environment variables")
}

func TestNewEnQuery(t *testing.T) {
	result := NewEnv()

	assert.Equal(t, "test1", (*result.envs)["TEST1"], "Must extract and parse environment variables")
	assert.Contains(t, "=test2=", (*result.envs)["TEST2"], "Must extract and parse environment variables")
}

func TestGetAllValues(t *testing.T) {
	setTestingEnvs()

	q := NewEnv()

	keys := q.GetAllValues()

	results := []string{}

	for _, v := range keys {
		if v == "test1" || v == "=test2=" {
			results = append(results, v)
		}
	}

	assert.Len(t, results, 2, "Must contains 2 elements")
}

func TestGetAllKeys(t *testing.T) {
	setTestingEnvs()

	q := NewEnv()

	keys := q.GetAllKeys()

	results := []string{}

	for _, k := range keys {
		if k == "TEST1" || k == "TEST2" {
			results = append(results, k)
		}
	}

	assert.Len(t, results, 2, "Must contains 2 elements")
}

func TestFindEntries(t *testing.T) {
	setTestingEnvs()

	q := NewEnv()

	keys, err := q.FindEntries(".*?1")

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, keys, 1, "Must contains 1 elements")
	assert.Equal(t, "test1", keys["TEST1"], "Must have env key and value")

	_, err = q.FindEntries("?")

	assert.EqualError(t, err, "error parsing regexp: missing argument to repetition operator: `?`", "Must return an error when regexp is unvalid")
}

func TestGetString(t *testing.T) {
	setTestingEnvs()

	q := NewEnv()

	value, err := q.GetString("TEST1")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "test1", value, "Must return value")

	value, err = q.GetString("TEST100")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, "", value, "Must return empty string")
}

func TestGetInt(t *testing.T) {
	err := os.Setenv("TEST3", "1")

	if err != nil {
		logrus.Fatal(err)
	}

	q := NewEnv()

	value, err := q.GetInt("TEST3")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, 1, value, "Must return value")

	value, err = q.GetInt("TEST100")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, 0, value, "Must return value")

	value, err = q.GetInt("TEST1")

	assert.EqualError(t, err, "Variable can't be converted", "Must return an error when variable can't be found")
	assert.Equal(t, 0, value, "Must return empty string")
}

func TestGetBool(t *testing.T) {
	err := os.Setenv("TEST4", "true")

	if err != nil {
		logrus.Fatal(err)
	}

	q := NewEnv()

	value, err := q.GetBool("TEST4")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, true, value, "Must return value")

	value, err = q.GetBool("TEST100")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, false, value, "Must return value")

	value, err = q.GetBool("TEST1")

	assert.EqualError(t, err, "Variable can't be converted", "Must return an error when variable can't be found")
	assert.Equal(t, false, value, "Must return empty string")
}

func TestGetFloat(t *testing.T) {
	err := os.Setenv("TEST5", "0.01")

	if err != nil {
		logrus.Fatal(err)
	}

	q := NewEnv()

	value, err := q.GetFloat("TEST5")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, float32(0.01), value, "Must return value")

	value, err = q.GetFloat("TEST100")

	assert.EqualError(t, err, "Variable not found", "Must return an error when variable can't be found")
	assert.Equal(t, float32(0), value, "Must return value")

	value, err = q.GetFloat("TEST1")

	assert.EqualError(t, err, "Variable can't be converted", "Must return an error when variable can't be found")
	assert.Equal(t, float32(0), value, "Must return empty string")
}

func TestCreateANode(t *testing.T) {
	n := newNode()

	assert.Equal(t, *n, node{childs: []*node{}}, "Must creates a new node")

	rootNode := newRootNode()

	assert.Equal(t, *rootNode, node{childs: []*node{}, root: true}, "Must creates a new root node")
}

func TestFindChildWithKey(t *testing.T) {
	root := newRootNode()

	node := newNode()
	node.key = "test"
	node.value = "value"
	root.appendChild(node)

	result, exists := root.findChildWithKey("test")

	assert.True(t, exists, "Must return true cause element was found")
	assert.Equal(t, node, result, "Must return child node with key test")

	_, exists = root.findChildWithKey("test1")

	assert.False(t, exists, "Must return false cause element was not found")
}

func TestAppendChild(t *testing.T) {
	root := newRootNode()

	node := newNode()
	node.key = "test"
	node.value = "value"

	result := root.appendChild(node)

	assert.True(t, result, "Must return true cause element was successfully added")
	assert.Equal(t, node, root.childs[0], "Must have node added as child")

	node2 := newNode()
	node2.key = "test"
	node2.value = "value2"

	result = root.appendChild(node2)

	assert.False(t, result, "Must return false cause an element with this key already exists")
	assert.Len(t, root.childs, 1, "Must still have one node")
	assert.Equal(t, node, root.childs[0], "Must have node added before")
}

func TestFindAllChildsWithKey(t *testing.T) {
	nodes := map[string]*node{}

	root := newRootNode()
	n := root

	var accumulatedKey string

	for _, i := range []string{"1", "2", "3"} {
		t := newNode()
		t.key = "test" + i
		t.value = "value" + i
		n.appendChild(t)

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
		n.appendChild(t)

		n = t

		if len(accumulatedKey) == 0 {
			accumulatedKey = i
		} else {
			accumulatedKey += "." + i
		}

		nodes[accumulatedKey] = t
	}

	results := root.findAllChildsWithKey("test3")

	assert.Equal(t, []*node{nodes["4.5.6.3"], nodes["1.2.3"]}, *results, "Must recurse over tree to find keys")
}

func TestAppendChildToTree(t *testing.T) {
	nodes := map[string]*node{}

	root := newRootNode()
	n := root

	var accumulatedKey string

	for _, i := range []string{"1", "2", "3"} {
		t := newNode()
		t.key = "test" + i
		t.value = "value" + i
		n.appendChild(t)

		n = t

		if len(accumulatedKey) == 0 {
			accumulatedKey = i
		} else {
			accumulatedKey += "." + i
		}

		nodes[accumulatedKey] = t
	}

	node4 := newNode()
	node4.key = "test4"
	node4.value = "value4"

	assert.True(t, root.appendChildToTree(node4, []string{"test1", "test2", "test3"}), "Must return true as added element must be successful")

	assert.Len(t, nodes["1.2.3"].childs, 1, "Must contains 1 element")
	assert.Equal(t, node4, nodes["1.2.3"].childs[0], "Must append element to node tree")

	node5 := newNode()
	node5.key = "test5"
	node5.value = "value5"

	assert.True(t, root.appendChildToTree(node5, []string{"test1", "test2", "test3"}), "Must return true as added element must be successful")

	assert.Len(t, nodes["1.2.3"].childs, 2, "Must contains 2 element")
	assert.Equal(t, node4, nodes["1.2.3"].childs[0], "Must recurse over tree to find keys")
	assert.Equal(t, node5, nodes["1.2.3"].childs[1], "Must append element to node tree")

	assert.False(t, root.appendChildToTree(node4, []string{"test1", "test2", "test3"}), "Must return false cause an element with this key already exists")

	assert.Len(t, nodes["1.2.3"].childs, 2, "Must contains 2 element")
}
