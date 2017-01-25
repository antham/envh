package envh

import (
	"regexp"
)

// Env manages environment variables
// by giving a convenient helper
// to interact with them
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

// GetString returns a string if variable exists
// or an error otherwise
func (e Env) GetString(key string) (string, error) {
	return getString(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
}

// GetInt returns an integer if variable exists
// or an error if value is not an integer or doesn't exist
func (e Env) GetInt(key string) (int, error) {
	return getInt(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
}

// GetFloat returns a float if variable exists
// or an error if value is not a float or doesn't exist
func (e Env) GetFloat(key string) (float32, error) {
	return getFloat(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
}

// GetBool returns a boolean if variable exists
// or an error if value is not a boolean or doesn't exist
func (e Env) GetBool(key string) (bool, error) {
	return getBool(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
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
