package envh

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func ExampleEnvTree_FindString() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.FindString("ENVH", "DB", "USERNAME"))
	fmt.Println(env.FindString("ENVH", "DB", "WHATEVER"))
	// Output:
	// foo <nil>
	//  Variable not found
}

func ExampleEnvTree_FindStringUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.FindStringUnsecured("ENVH", "DB", "USERNAME"))
	fmt.Println(env.FindStringUnsecured("ENVH", "DB", "WHATEVER"))
	// Output:
	// foo
	//
}

func ExampleEnvTree_FindInt() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.FindInt("ENVH", "DB", "PORT"))
	fmt.Println(env.FindInt("ENVH", "DB", "USERNAME"))
	fmt.Println(env.FindInt("ENVH", "DB", "WHATEVER"))
	// Output:
	// 3306 <nil>
	// 0 Value "foo" can't be converted to type "int"
	// 0 Variable not found
}

func ExampleEnvTree_FindIntUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.FindIntUnsecured("ENVH", "DB", "PORT"))
	fmt.Println(env.FindIntUnsecured("ENVH", "DB", "USERNAME"))
	fmt.Println(env.FindIntUnsecured("ENVH", "DB", "WHATEVER"))
	// Output:
	// 3306
	// 0
	// 0
}

func ExampleEnvTree_FindBool() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.FindBool("ENVH", "MAILER", "ENABLED"))
	fmt.Println(env.FindBool("ENVH", "DB", "USERNAME"))
	fmt.Println(env.FindBool("ENVH", "DB", "WHATEVER"))
	// Output:
	// true <nil>
	// false Value "foo" can't be converted to type "bool"
	// false Variable not found
}

func ExampleEnvTree_FindBoolUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.FindBoolUnsecured("ENVH", "MAILER", "ENABLED"))
	fmt.Println(env.FindBoolUnsecured("ENVH", "DB", "USERNAME"))
	fmt.Println(env.FindBoolUnsecured("ENVH", "DB", "WHATEVER"))
	// Output:
	// true
	// false
	// false
}

func ExampleEnvTree_FindFloat() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.FindFloat("ENVH", "DB", "USAGE", "LIMIT"))
	fmt.Println(env.FindFloat("ENVH", "DB", "USERNAME"))
	fmt.Println(env.FindFloat("ENVH", "DB", "WHATEVER"))
	// Output:
	// 95.6 <nil>
	// 0 Value "foo" can't be converted to type "float"
	// 0 Variable not found
}

func ExampleEnvTree_FindFloatUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.FindFloatUnsecured("ENVH", "DB", "USAGE", "LIMIT"))
	fmt.Println(env.FindFloatUnsecured("ENVH", "DB", "USERNAME"))
	fmt.Println(env.FindFloatUnsecured("ENVH", "DB", "WHATEVER"))
	// Output:
	// 95.6
	// 0
	// 0
}

func ExampleEnvTree_IsExistingSubTree() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.IsExistingSubTree("ENVH", "MAILER", "HOST"))
	fmt.Println(env.IsExistingSubTree("ENVH", "MAILER", "WHATEVER"))
	// Output:
	// true
	// false
}

func ExampleEnvTree_HasSubTreeValue() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.HasSubTreeValue("ENVH", "MAILER", "HOST"))
	fmt.Println(env.HasSubTreeValue("ENVH", "MAILER"))
	fmt.Println(env.HasSubTreeValue("ENVH", "MAILER", "WHATEVER"))
	// Output:
	// true <nil>
	// false <nil>
	// false No node found at path "ENVH -> MAILER -> WHATEVER"
}

func ExampleEnvTree_HasSubTreeValueUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	fmt.Println(env.HasSubTreeValueUnsecured("ENVH", "MAILER", "HOST"))
	fmt.Println(env.HasSubTreeValueUnsecured("ENVH", "MAILER"))
	fmt.Println(env.HasSubTreeValueUnsecured("ENVH", "MAILER", "WHATEVER"))
	// Output:
	// true
	// false
	// false
}

func ExampleEnvTree_HasValue() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	usernameTree, err := env.FindSubTree("ENVH", "DB", "USERNAME")

	if err != nil {
		return
	}

	fmt.Println(env.HasValue())
	fmt.Println(usernameTree.HasValue())
	// Output:
	// false
	// true
}

func ExampleEnvTree_FindSubTree() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	dbTree, err := env.FindSubTree("ENVH", "DB")
	dbChildrenKeys := dbTree.GetChildrenKeys()
	sort.Strings(dbChildrenKeys)

	fmt.Print(dbChildrenKeys)
	fmt.Print(" ")
	fmt.Println(err)

	mailerTree, err := env.FindSubTree("ENVH", "MAILER")
	mailerChildrenKeys := mailerTree.GetChildrenKeys()
	sort.Strings(mailerChildrenKeys)

	fmt.Print(mailerChildrenKeys)
	fmt.Print(" ")
	fmt.Println(err)

	fmt.Println(env.FindSubTree("ENVH", "MAILER", "WHATEVER"))
	// Output:
	// [PASSWORD PORT USAGE USERNAME] <nil>
	// [ENABLED HOST PASSWORD USERNAME] <nil>
	// {<nil>} No node found at path "ENVH -> MAILER -> WHATEVER"
}

func ExampleEnvTree_FindSubTreeUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	dbTree := env.FindSubTreeUnsecured("ENVH", "DB")
	dbChildrenKeys := dbTree.GetChildrenKeys()
	sort.Strings(dbChildrenKeys)

	fmt.Println(dbChildrenKeys)

	mailerTree := env.FindSubTreeUnsecured("ENVH", "MAILER")
	mailerChildrenKeys := mailerTree.GetChildrenKeys()
	sort.Strings(mailerChildrenKeys)

	fmt.Println(mailerChildrenKeys)

	fmt.Println(env.FindSubTreeUnsecured("ENVH", "MAILER", "WHATEVER"))
	// Output:
	// [PASSWORD PORT USAGE USERNAME]
	// [ENABLED HOST PASSWORD USERNAME]
	// {<nil>}
}

func ExampleEnvTree_GetKey() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	dbTree, err := env.FindSubTree("ENVH", "DB")

	if err != nil {
		return
	}

	fmt.Println(dbTree.GetKey())
	// Output: DB
}

func ExampleEnvTree_FindChildrenKeys() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	children, err := env.FindChildrenKeys("ENVH", "DB")

	sort.Strings(children)

	fmt.Print(children)
	fmt.Print(" ")
	fmt.Println(err)
	fmt.Println(env.FindChildrenKeys("ENVH", "WHATEVER"))
	// Output:
	// [PASSWORD PORT USAGE USERNAME] <nil>
	// [] No node found at path "ENVH -> WHATEVER"
}

func ExampleEnvTree_FindChildrenKeysUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	children := env.FindChildrenKeysUnsecured("ENVH", "DB")

	sort.Strings(children)

	fmt.Println(children)
	fmt.Println(env.FindChildrenKeysUnsecured("ENVH", "WHATEVER"))
	// Output:
	// [PASSWORD PORT USAGE USERNAME]
	// []
}

func ExampleEnvTree_GetBool() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	enabledTree, err := env.FindSubTree("ENVH", "MAILER", "ENABLED")

	if err != nil {
		return
	}

	fmt.Println(env.GetBool())
	fmt.Println(enabledTree.GetBool())
	// Output:
	// false Variable not found
	// true <nil>
}

func ExampleEnvTree_GetBoolUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	enabledTree, err := env.FindSubTree("ENVH", "MAILER", "ENABLED")

	if err != nil {
		return
	}

	fmt.Println(env.GetBoolUnsecured())
	fmt.Println(enabledTree.GetBoolUnsecured())
	// Output:
	// false
	// true
}

func ExampleEnvTree_GetInt() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	portTree, err := env.FindSubTree("ENVH", "DB", "PORT")

	if err != nil {
		return
	}

	fmt.Println(env.GetInt())
	fmt.Println(portTree.GetInt())
	// Output:
	// 0 Variable not found
	// 3306 <nil>
}

func ExampleEnvTree_GetIntUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	portTree, err := env.FindSubTree("ENVH", "DB", "PORT")

	if err != nil {
		return
	}

	fmt.Println(env.GetIntUnsecured())
	fmt.Println(portTree.GetIntUnsecured())
	// Output:
	// 0
	// 3306
}

func ExampleEnvTree_GetString() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	usernameTree, err := env.FindSubTree("ENVH", "DB", "USERNAME")

	if err != nil {
		return
	}

	fmt.Println(env.GetString())
	fmt.Println(usernameTree.GetString())
	// Output:
	//  Variable not found
	// foo <nil>
}

func ExampleEnvTree_GetStringUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	usernameTree, err := env.FindSubTree("ENVH", "DB", "USERNAME")

	if err != nil {
		return
	}

	fmt.Println(env.GetStringUnsecured())
	fmt.Println(usernameTree.GetStringUnsecured())
	// Output:
	//
	// foo
}

func ExampleEnvTree_GetFloat() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	portTree, err := env.FindSubTree("ENVH", "DB", "USAGE", "LIMIT")

	if err != nil {
		return
	}

	fmt.Println(env.GetFloat())
	fmt.Println(portTree.GetFloat())
	// Output:
	// 0 Variable not found
	// 95.6 <nil>
}

func ExampleEnvTree_GetFloatUnsecured() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	portTree, err := env.FindSubTree("ENVH", "DB", "USAGE", "LIMIT")

	if err != nil {
		return
	}

	fmt.Println(env.GetFloatUnsecured())
	fmt.Println(portTree.GetFloatUnsecured())
	// Output:
	// 0
	// 95.6
}

func ExampleEnvTree_GetChildrenKeys() {
	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGE_LIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	dbTree, err := env.FindSubTree("ENVH", "DB")

	if err != nil {
		return
	}

	children := dbTree.GetChildrenKeys()

	sort.Strings(children)

	fmt.Println(children)
	// Output:
	// [PASSWORD PORT USAGE USERNAME]
}

func ExampleEnvTree_PopulateStruct() {
	type ENVH struct {
		DB struct {
			USERNAME   string
			PASSWORD   string
			PORT       int
			USAGELIMIT float32
		}
		MAILER struct {
			HOST     string
			USERNAME string
			PASSWORD string
			ENABLED  bool
		}
	}

	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGELIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")
	setEnv("ENVH_MAILER_ENABLED", "true")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	s := ENVH{}

	err = env.PopulateStruct(&s)

	if err != nil {
		return
	}

	b, err := json.Marshal(s)

	if err != nil {
		return
	}

	fmt.Println(string(b))
	// Output:
	// {"DB":{"USERNAME":"foo","PASSWORD":"bar","PORT":3306,"USAGELIMIT":95.6},"MAILER":{"HOST":"127.0.0.1","USERNAME":"foo","PASSWORD":"bar","ENABLED":true}}
}

func ExampleEnvTree_PopulateStructWithStrictMode() {
	type ENVH struct {
		DB struct {
			USERNAME   string
			PASSWORD   string
			PORT       int
			USAGELIMIT float32
		}
		MAILER struct {
			HOST     string
			USERNAME string
			PASSWORD string
			ENABLED  bool
		}
	}

	os.Clearenv()
	setEnv("ENVH_DB_USERNAME", "foo")
	setEnv("ENVH_DB_PASSWORD", "bar")
	setEnv("ENVH_DB_PORT", "3306")
	setEnv("ENVH_DB_USAGELIMIT", "95.6")
	setEnv("ENVH_MAILER_HOST", "127.0.0.1")
	setEnv("ENVH_MAILER_USERNAME", "foo")
	setEnv("ENVH_MAILER_PASSWORD", "bar")

	env, err := NewEnvTree("^ENVH", "_")

	if err != nil {
		return
	}

	s := ENVH{}

	err = env.PopulateStructWithStrictMode(&s)

	fmt.Println(err)
	// Output:
	// Variable not found
}
