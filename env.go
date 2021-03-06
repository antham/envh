// Package envh provides convenient helpers to manage easily your environment variables.
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

// GetStringUnsecured is insecured version of GetString to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing, it returns default zero string value.
// This function has to be used carefully
func (e Env) GetStringUnsecured(key string) string {
	if val, err := getString(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	}); err == nil {
		return val
	}

	return ""
}

// GetInt returns an integer if variable exists
// or an error if value is not an integer or doesn't exist
func (e Env) GetInt(key string) (int, error) {
	return getInt(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
}

// GetIntUnsecured is insecured version of GetInt to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing or not an int value, it returns default zero int value.
// This function has to be used carefully
func (e Env) GetIntUnsecured(key string) int {
	if val, err := getInt(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	}); err == nil {
		return val
	}

	return 0
}

// GetFloat returns a float if variable exists
// or an error if value is not a float or doesn't exist
func (e Env) GetFloat(key string) (float32, error) {
	return getFloat(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
}

// GetFloatUnsecured is insecured version of GetFloat to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing or not a floating value, it returns default zero floating value.
// This function has to be used carefully
func (e Env) GetFloatUnsecured(key string) float32 {
	if val, err := getFloat(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	}); err == nil {
		return val
	}

	return 0
}

// GetBool returns a boolean if variable exists
// or an error if value is not a boolean or doesn't exist
func (e Env) GetBool(key string) (bool, error) {
	return getBool(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	})
}

// GetBoolUnsecured is insecured version of GetBool to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing or not a boolean value, it returns default zero boolean value.
// This function has to be used carefully
func (e Env) GetBoolUnsecured(key string) bool {
	if val, err := getBool(func() (string, bool) {
		v, ok := (*e.envs)[key]

		return v, ok
	}); err == nil {
		return val
	}

	return false
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

// FindEntriesUnsecured is insecured version of FindEntriesUnsecured to avoid the burden
// of rechecking errors if it was done already. If any errors occurred cause
// the variable is missing or not a boolean value, it returns default empty map.
// This function has to be used carefully.
func (e Env) FindEntriesUnsecured(reg string) map[string]string {
	results := map[string]string{}

	r, err := regexp.Compile(reg)

	if err != nil {
		return results
	}

	for k, v := range *e.envs {
		if r.MatchString(k) {
			results[k] = v
		}
	}

	return results
}
