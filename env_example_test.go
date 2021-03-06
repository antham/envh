package envh

import (
	"fmt"
	"os"
	"sort"
)

func ExampleEnv_GetAllKeys() {
	os.Clearenv()
	setEnv("HELLO", "world")
	setEnv("FOO", "bar")

	env := NewEnv()

	keys := env.GetAllKeys()

	sort.Strings(keys)

	fmt.Println(keys)
	// Output: [FOO HELLO]
}

func ExampleEnv_GetAllValues() {
	os.Clearenv()
	setEnv("HELLO", "world")
	setEnv("FOO", "bar")

	env := NewEnv()

	values := env.GetAllValues()

	sort.Strings(values)

	fmt.Println(values)
	// Output: [bar world]
}

func ExampleEnv_GetString() {
	os.Clearenv()
	setEnv("HELLO", "world")

	env := NewEnv()

	fmt.Println(env.GetString("HELLO"))
	// Output: world <nil>
}

func ExampleEnv_GetStringUnsecured() {
	os.Clearenv()
	setEnv("HELLO", "world")

	env := NewEnv()

	fmt.Println(env.GetStringUnsecured("HELLO"))
	// Output: world
}

func ExampleEnv_GetInt() {
	os.Clearenv()
	setEnv("INT", "1")
	setEnv("STRING", "TEST")

	env := NewEnv()

	fmt.Println(env.GetInt("INT"))
	fmt.Println(env.GetInt("STRING"))

	// Output:
	// 1 <nil>
	// 0 Value "TEST" can't be converted to type "int"
}

func ExampleEnv_GetIntUnsecured() {
	os.Clearenv()
	setEnv("INT", "1")
	setEnv("STRING", "TEST")

	env := NewEnv()

	fmt.Println(env.GetIntUnsecured("INT"))
	fmt.Println(env.GetIntUnsecured("STRING"))

	// Output:
	// 1
	// 0
}

func ExampleEnv_GetFloat() {
	os.Clearenv()
	setEnv("FLOAT", "1.1")
	setEnv("STRING", "TEST")

	env := NewEnv()

	f, err := env.GetFloat("FLOAT")

	fmt.Printf("%0.1f ", f)
	fmt.Println(err)
	fmt.Println(env.GetFloat("STRING"))

	// Output:
	// 1.1 <nil>
	// 0 Value "TEST" can't be converted to type "float"
}

func ExampleEnv_GetFloatUnsecured() {
	os.Clearenv()
	setEnv("FLOAT", "1.1")
	setEnv("STRING", "TEST")

	env := NewEnv()

	fmt.Printf("%0.1f\n", env.GetFloatUnsecured("FLOAT"))
	fmt.Println(env.GetFloatUnsecured("STRING"))

	// Output:
	// 1.1
	// 0
}

func ExampleEnv_GetBool() {
	os.Clearenv()
	setEnv("BOOL", "true")
	setEnv("STRING", "TEST")

	env := NewEnv()

	fmt.Println(env.GetBool("BOOL"))
	fmt.Println(env.GetBool("STRING"))

	// Output:
	// true <nil>
	// false Value "TEST" can't be converted to type "bool"
}

func ExampleEnv_GetBoolUnsecured() {
	os.Clearenv()
	setEnv("BOOL", "true")
	setEnv("STRING", "TEST")

	env := NewEnv()

	fmt.Println(env.GetBoolUnsecured("BOOL"))
	fmt.Println(env.GetBoolUnsecured("STRING"))

	// Output:
	// true
	// false
}

func ExampleEnv_FindEntries() {
	os.Clearenv()
	setEnv("API_USERNAME", "user")
	setEnv("API_PASSWORD", "password")
	setEnv("DB_USERNAME", "user")
	setEnv("DB_PASSWORD", "user")

	env := NewEnv()

	entries, err := env.FindEntries("API.*")

	fmt.Printf("API -> PASSWORD = %s, API -> USERNAME = %s ", entries["API_PASSWORD"], entries["API_PASSWORD"])
	fmt.Println(err)
	fmt.Println(env.FindEntries("*"))

	// Output:
	// API -> PASSWORD = password, API -> USERNAME = password <nil>
	// map[] error parsing regexp: missing argument to repetition operator: `*`
}

func ExampleEnv_FindEntriesUnsecured() {
	os.Clearenv()
	setEnv("API_USERNAME", "user")
	setEnv("API_PASSWORD", "password")
	setEnv("DB_USERNAME", "user")
	setEnv("DB_PASSWORD", "user")

	env := NewEnv()

	entries := env.FindEntriesUnsecured("API.*")

	fmt.Printf("API -> PASSWORD = %s, API -> USERNAME = %s\n", entries["API_PASSWORD"], entries["API_PASSWORD"])
	fmt.Println(env.FindEntriesUnsecured("*"))

	// Output:
	// API -> PASSWORD = password, API -> USERNAME = password
	// map[]
}
