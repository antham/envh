package envh

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

// Env manage environment variables
// by giving high level api to interact with them
type Env struct {
	envs *map[string]string
}

// NewEnv creates a new Env instance
func NewEnv() Env {
	return Env{parseVars()}
}

// GetAllValues retrieves a slice of all environment variables values
func (e Env) GetAllValues() []string {
	results := []string{}

	for _, v := range *e.envs {
		results = append(results, v)
	}

	return results
}

// GetAllKeys retrieves a slice of all environment variables keys
func (e Env) GetAllKeys() []string {
	results := []string{}

	for k := range *e.envs {
		results = append(results, k)
	}

	return results
}

// FindEntries retrieves all keys matching a given regexp and their
// corresponding values
func (e Env) FindEntries(reg string) (map[string]string, error) {
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
