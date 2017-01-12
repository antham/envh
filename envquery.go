package envquery

import (
	"os"
	"strings"
)

func parseVars() *map[string]string {
	results := map[string]string{}

	for _, v := range os.Environ() {
		e := strings.SplitN(v, "=", 2)

		results[e[0]] = e[1]
	}

	return &results
}

// EnvQuery manage environment variables
// by giving high level api to interact with them
type EnvQuery struct {
	envs *map[string]string
}

// NewEnvQuery creates a new EnvQuery instance
func NewEnvQuery() EnvQuery {
	return EnvQuery{parseVars()}
}

