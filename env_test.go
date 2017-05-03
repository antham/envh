package envh

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewEnvQuery(t *testing.T) {
	setTestingEnvs()
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

func TestFindEntriesUnsecured(t *testing.T) {
	setTestingEnvs()

	q := NewEnv()

	keys := q.FindEntriesUnsecured(".*?1")

	assert.Len(t, keys, 1, "Must contains 1 elements")
	assert.Equal(t, "test1", keys["TEST1"], "Must have env key and value")

	keys = q.FindEntriesUnsecured("?")
	assert.Len(t, keys, 0, "Must contains 0 elements")
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

func TestGetStringUnsecured(t *testing.T) {
	setTestingEnvs()

	q := NewEnv()

	value := q.GetStringUnsecured("TEST1")

	assert.Equal(t, "test1", value, "Must return value")

	value = q.GetStringUnsecured("TEST100")

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

	assert.EqualError(t, err, `Value "test1" can't be converted to type "int"`, "Must return an error when variable can't be found")
	assert.Equal(t, 0, value, "Must return empty string")
}

func TestGetIntUnsecured(t *testing.T) {
	err := os.Setenv("TEST3", "1")

	if err != nil {
		logrus.Fatal(err)
	}

	q := NewEnv()

	value := q.GetIntUnsecured("TEST3")

	assert.Equal(t, 1, value, "Must return value")

	value = q.GetIntUnsecured("TEST100")

	assert.Equal(t, 0, value, "Must return value")

	value = q.GetIntUnsecured("TEST1")

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

	assert.EqualError(t, err, `Value "test1" can't be converted to type "bool"`, "Must return an error when variable can't be found")
	assert.Equal(t, false, value, "Must return empty string")
}

func TestGetBoolUnsecured(t *testing.T) {
	err := os.Setenv("TEST4", "true")

	if err != nil {
		logrus.Fatal(err)
	}

	q := NewEnv()

	value := q.GetBoolUnsecured("TEST4")

	assert.Equal(t, true, value, "Must return value")

	value = q.GetBoolUnsecured("TEST100")

	assert.Equal(t, false, value, "Must return value")

	value = q.GetBoolUnsecured("TEST1")

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

	assert.EqualError(t, err, `Value "test1" can't be converted to type "float"`, "Must return an error when variable can't be found")
	assert.Equal(t, float32(0), value, "Must return empty string")
}

func TestGetFloatUnsecured(t *testing.T) {
	err := os.Setenv("TEST5", "0.01")

	if err != nil {
		logrus.Fatal(err)
	}

	q := NewEnv()

	value := q.GetFloatUnsecured("TEST5")

	assert.Equal(t, float32(0.01), value, "Must return value")

	value = q.GetFloatUnsecured("TEST100")

	assert.Equal(t, float32(0), value, "Must return value")

	value = q.GetFloatUnsecured("TEST1")

	assert.Equal(t, float32(0), value, "Must return empty string")
}
