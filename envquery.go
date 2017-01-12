package envquery

import (
	"os"
	"regexp"
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

// GetAllValues retrieves a slice of all environment variables values
func (e EnvQuery) GetAllValues() []string {
	results := []string{}

	for _, v := range *e.envs {
		results = append(results, v)
	}

	return results
}

// GetAllKeys retrieves a slice of all environment variables keys
func (e EnvQuery) GetAllKeys() []string {
	results := []string{}

	for k := range *e.envs {
		results = append(results, k)
	}

	return results
}

// FindEntries retrieves all keys matching a given regexp and their
// corresponding values
func (e EnvQuery) FindEntries(reg string) (map[string]string, error) {
	results := map[string]string{}

	r, err := regexp.Compile(reg)

	if err != nil {
		return results, err
	}

	for k, v := range *e.envs {
		if r.MatchString(k) {
			results[k] = v
		}
	}

	return results, nil
}
