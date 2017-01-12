package envquery

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVars(t *testing.T) {
	datas := map[string]string{
		"TEST1": "test",
		"TEST2": "=test=",
	}

	for k, v := range datas {
		os.Setenv(k, v)
	}

	result := parseVars()

	assert.Equal(t, "test", (*result)["TEST1"], "Must extract and parse environment variables")
	assert.Contains(t, "=test=", (*result)["TEST2"], "Must extract and parse environment variables")
}
