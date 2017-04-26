package envh

import (
	"fmt"
	"strings"

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

type SUM struct {
	LEFTOPERAND  int
	RIGHTOPERAND int
	RESULT       int
}

func (s *SUM) Walk(tree *EnvTree, keyChain []string) (bool, error) {
	if iterator, ok := map[string]func(*EnvTree, []string) (bool, error){
		"SUM_LEFTOPERAND": s.validateLeftOperand,
		"SUM_RESULT":      s.setResult,
	}[strings.Join(keyChain, "_")]; ok {
		return iterator(tree, keyChain)
	}

	return false, nil

}

func (s *SUM) setResult(tree *EnvTree, keyChain []string) (bool, error) {
	left, err := tree.FindInt("SUM", "LEFTOPERAND")

	if err != nil {
		return true, fmt.Errorf(`Can't find "SUM_LEFTOPERAND"`)
	}

	right, err := tree.FindInt("SUM", "RIGHTOPERAND")

	if err != nil {
		return true, fmt.Errorf(`Can't find "SUM_LEFTOPERAND"`)
	}

	s.RESULT = left + right

	return true, nil
}

func (s *SUM) validateLeftOperand(tree *EnvTree, keyChain []string) (bool, error) {
	val, err := tree.FindInt(keyChain...)

	if err != nil {
		return false, fmt.Errorf(`"SUM_LEFTOPERAND" can't be found`)
	}

	if val <= 0 {
		return false, fmt.Errorf(`"LEFTOPERAND" must be greater than 0`)
	}

	return false, nil
}

func TestPopulateStructWithCustomSet(t *testing.T) {
	setEnv("SUM_LEFTOPERAND", "1")
	setEnv("SUM_RIGHTOPERAND", "2")

	actual := SUM{}

	tree, err := NewEnvTree("SUM", "_")

	assert.NoError(t, err)

	err = populateStructFromEnvTree(&actual, &tree, false)

	assert.NoError(t, err)

	expected := SUM{LEFTOPERAND: 1, RIGHTOPERAND: 2, RESULT: 3}

	assert.Equal(t, expected, actual, "Must set result field")

	restoreEnvs()
}

func TestPopulateStructWithCustomSetTriggeringAnError(t *testing.T) {
	setEnv("SUM_LEFTOPERAND", "2")

	actual := SUM{}

	tree, err := NewEnvTree("SUM", "_")

	assert.NoError(t, err)

	err = populateStructFromEnvTree(&actual, &tree, false)

	assert.EqualError(t, err, `Can't find "SUM_LEFTOPERAND"`, "Must bubble up an error from Populate function")

	restoreEnvs()
}

func TestPopulateStructWithCustomValidationTriggeringAnError(t *testing.T) {
	setEnv("SUM_LEFTOPERAND", "0")
	setEnv("SUM_RIGHTOPERAND", "2")

	actual := SUM{}

	tree, err := NewEnvTree("SUM", "_")

	assert.NoError(t, err)

	err = populateStructFromEnvTree(&actual, &tree, false)

	assert.EqualError(t, err, `"LEFTOPERAND" must be greater than 0`, "Must validate data")

	restoreEnvs()
}
