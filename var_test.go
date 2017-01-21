package envh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVars(t *testing.T) {
	setTestingEnvs()
	result := parseVars()

	assert.Equal(t, "test1", (*result)["TEST1"], "Must extract and parse environment variables")
	assert.Contains(t, "=test2=", (*result)["TEST2"], "Must extract and parse environment variables")
}
