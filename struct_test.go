package envh

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPopulateStructWithNoPointersGiven(t *testing.T) {
	type POPULATESTRUCT struct{}

	tree, err := NewEnvTree("test", "_")

	assert.NoError(t, err, "Must return no errors")

	err = populateStructFromEnvTree(POPULATESTRUCT{}, &tree, false)

	assert.EqualError(t, err, `Type "struct" is not supported : you must provide "pointer to struct"`)

	err = populateStructFromEnvTree(8, &tree, false)

	assert.EqualError(t, err, `Type "int" is not supported : you must provide "pointer to struct"`)
}

func TestPopulateStructWithInnerPointer(t *testing.T) {
	type TEST7 struct {
		TEST8  int
		TEST9  float32
		TEST10 string
		TEST11 bool
	}

	type TEST4 struct {
		TEST5 *string
		TEST6 TEST7
	}

	type POPULATESTRUCT struct {
		TEST2 string
		TEST3 TEST4
	}

	actual := POPULATESTRUCT{}

	tree, err := NewEnvTree("test", "_")

	assert.NoError(t, err)

	err = populateStructFromEnvTree(&actual, &tree, false)

	restoreEnvs()

	assert.EqualError(t, err, `Type "ptr" is not supported : you must provide "int32, float32, string, boolean or struct"`)
}

func TestPopulateStructWithTypeErrors(t *testing.T) {
	type TEST5 struct {
		TEST6 int
	}

	type POPULATESTRUCT struct {
		TEST1 float32
		TEST2 int
		TEST3 bool
		TEST4 TEST5
	}

	type g struct {
		init       func()
		checkError func(err error)
		teardown   func()
	}

	tests := []g{
		{
			init: func() {
				setEnv("POPULATESTRUCT_TEST1", "value1")
			},
			checkError: func(err error) {
				assert.EqualError(t, err, `Value "value1" can't be converted to type "float"`)
			},
		},
		{
			init: func() {
				setEnv("POPULATESTRUCT_TEST2", "value2")
			},
			checkError: func(err error) {
				assert.EqualError(t, err, `Value "value2" can't be converted to type "int"`)
			},
		},
		{
			init: func() {
				setEnv("POPULATESTRUCT_TEST3", "value3")
			},
			checkError: func(err error) {
				assert.EqualError(t, err, `Value "value3" can't be converted to type "bool"`)
			},
		},
		{
			init: func() {
				setEnv("POPULATESTRUCT_TEST4_TEST6", "value4")
			},
			checkError: func(err error) {
				assert.EqualError(t, err, `Value "value4" can't be converted to type "int"`)
			},
		},
	}

	for _, s := range tests {
		actual := POPULATESTRUCT{}

		s.init()

		tree, err := NewEnvTree("POPULATESTRUCT", "_")

		assert.NoError(t, err)

		err = populateStructFromEnvTree(&actual, &tree, false)
		s.checkError(err)
		restoreEnvs()
	}
}

func TestPopulateStructWithStrictCheckDisabled(t *testing.T) {
	type TEST7 struct {
		TEST8  int
		TEST9  float32
		TEST10 string
		TEST11 bool
	}

	type TEST4 struct {
		TEST5 string
		TEST6 TEST7
	}

	type POPULATESTRUCT struct {
		TEST2 string
		TEST3 TEST4
	}

	actual := POPULATESTRUCT{}

	setEnv("POPULATESTRUCT_TEST2", "value1")
	setEnv("POPULATESTRUCT_TEST3_TEST6_TEST8", "1")
	setEnv("POPULATESTRUCT_TEST3_TEST6_TEST9", "1.1")
	setEnv("POPULATESTRUCT_TEST3_TEST6_TEST10", "value test 10")
	setEnv("POPULATESTRUCT_TEST3_TEST6_TEST11", "true")

	tree, err := NewEnvTree("POPULATESTRUCT", "_")

	assert.NoError(t, err)

	err = populateStructFromEnvTree(&actual, &tree, false)

	restoreEnvs()

	expected := POPULATESTRUCT{
		"value1",
		TEST4{
			"",
			TEST7{
				1,
				1.1,
				"value test 10",
				true,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual, "Must popuplate struct with value given by environment variables")
}

func TestPopulateStructWithStrictCheckEnabled(t *testing.T) {
	type POPULATESTRUCT struct {
		TEST8  int
		TEST9  float32
		TEST10 string
		TEST11 bool
	}

	type g struct {
		init       func()
		checkError func(err error)
		teardown   func()
	}

	tests := []g{
		{
			init: func() {
			},
			checkError: func(err error) {
				assert.EqualError(t, err, `Variable not found`)
			},
		},
		{
			init: func() {
				setEnv("POPULATESTRUCT_TEST8", "1")
			},
			checkError: func(err error) {
				assert.EqualError(t, err, `Variable not found`)
			},
		},
		{
			init: func() {
				setEnv("POPULATESTRUCT_TEST8", "1")
				setEnv("POPULATESTRUCT_TEST9", "1.1")
			},
			checkError: func(err error) {
				assert.EqualError(t, err, `Variable not found`)
			},
		},
		{
			init: func() {
				setEnv("POPULATESTRUCT_TEST8", "1")
				setEnv("POPULATESTRUCT_TEST9", "1.1")
				setEnv("POPULATESTRUCT_TEST10", "test")
			},
			checkError: func(err error) {
				assert.EqualError(t, err, `Variable not found`)
			},
		},
	}

	for _, s := range tests {
		actual := POPULATESTRUCT{}

		s.init()

		tree, err := NewEnvTree("POPULATESTRUCT", "_")

		assert.NoError(t, err)

		err = populateStructFromEnvTree(&actual, &tree, true)
		s.checkError(err)
		restoreEnvs()
	}
}
