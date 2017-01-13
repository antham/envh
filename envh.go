package envh

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ErrNotFound is triggered when environment variable cannot be found
var ErrNotFound = fmt.Errorf("Variable not found")

// ErrWrongType is triggered when we try to convert variable to a wrong type
var ErrWrongType = fmt.Errorf("Variable can't be converted")

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

// GetString return a string if variable exists
// or an error otherwise
func (e Env) GetString(key string) (string, error) {
	if v, ok := (*e.envs)[key]; ok {
		return v, nil
	}

	return "", ErrNotFound
}

// GetInt return an integer if variable exists
// or an error if value is not an integer or doesn't exist
func (e Env) GetInt(key string) (int, error) {
	v, ok := (*e.envs)[key]

	if !ok {
		return 0, ErrNotFound
	}

	i, err := strconv.Atoi(v)

	if err != nil {
		return 0, ErrWrongType
	}

	return i, nil
}

// GetFloat return a float if variable exists
// or an error if value is not a float or doesn't exist
func (e Env) GetFloat(key string) (float32, error) {
	v, ok := (*e.envs)[key]

	if !ok {
		return 0, ErrNotFound
	}

	f, err := strconv.ParseFloat(v, 32)

	if err != nil {
		return 0, ErrWrongType
	}

	return float32(f), nil
}

// GetBool return a boolean if variable exists
// or an error if value is not a boolean or doesn't exist
func (e Env) GetBool(key string) (bool, error) {
	v, ok := (*e.envs)[key]

	if !ok {
		return false, ErrNotFound
	}

	b, err := strconv.ParseBool(v)

	if err != nil {
		return false, ErrWrongType
	}

	return b, nil
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
