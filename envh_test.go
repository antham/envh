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
