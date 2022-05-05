# Envh [![codecov](https://codecov.io/gh/antham/envh/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/envh) [![Go Report Card](https://goreportcard.com/badge/github.com/antham/envh)](https://goreportcard.com/report/github.com/antham/envh) [![GoDoc](https://godoc.org/github.com/antham/envh?status.svg)](http://godoc.org/github.com/antham/envh) [![GitHub tag](https://img.shields.io/github/tag/antham/envh.svg)]()

This library is made up of two parts :

- Env object : it wraps your environments variables in an object and provides convenient helpers.
- Env tree object : it manages environment variables through a tree structure to store a config the same way as in a yaml file or whatever format allows to store a config hierarchically

## Install

    go get github.com/antham/envh

## How it works

Check [the godoc](http://godoc.org/github.com/antham/envh), there are many examples provided.

## Example with a tree dumped in a config struct

```go
package envh

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type CONFIG2 struct {
	DB struct {
		USERNAME   string
		PASSWORD   string
		HOST       string
		NAME       string
		PORT       int
		URL        string
		USAGELIMIT float32
	}
	MAILER struct {
		HOST     string
		USERNAME string
		PASSWORD string
		ENABLED  bool
	}
	MAP map[string]string
}

func (c *CONFIG2) Walk(tree *EnvTree, keyChain []string) (bool, error) {
	if setter, ok := map[string]func(*EnvTree, []string) error{
		"CONFIG2_DB_URL": c.setURL,
		"CONFIG2_MAP":    c.setMap,
	}[strings.Join(keyChain, "_")]; ok {
		return true, setter(tree, keyChain)
	}

	return false, nil
}

func (c *CONFIG2) setMap(tree *EnvTree, keyChain []string) error {
	datas := map[string]string{}

	keys, err := tree.FindChildrenKeys(keyChain...)

	if err != nil {
		return err
	}

	for _, key := range keys {
		value, err := tree.FindString(append(keyChain, key)...)

		if err != nil {
			return err
		}

		datas[key] = value
	}

	c.MAP = datas

	return nil
}

func (c *CONFIG2) setURL(tree *EnvTree, keyChain []string) error {
	datas := map[string]string{}

	for _, key := range []string{"USERNAME", "PASSWORD", "HOST", "NAME"} {
		value, err := tree.FindString("CONFIG2", "DB", key)

		if err != nil {
			return err
		}

		datas[key] = value
	}

	port, err := tree.FindInt("CONFIG2", "DB", "PORT")

	if err != nil {
		return err
	}

	c.DB.URL = fmt.Sprintf("jdbc:mysql://%s:%d/%s?user=%s&password=%s", datas["HOST"], port, datas["NAME"], datas["USERNAME"], datas["PASSWORD"])

	return nil
}

func ExampleStructWalker_customFieldSet() {
	os.Clearenv()
	setEnv("CONFIG2_DB_USERNAME", "foo")
	setEnv("CONFIG2_DB_PASSWORD", "bar")
	setEnv("CONFIG2_DB_HOST", "localhost")
	setEnv("CONFIG2_DB_NAME", "my-db")
	setEnv("CONFIG2_DB_PORT", "3306")
	setEnv("CONFIG2_DB_USAGELIMIT", "95.6")
	setEnv("CONFIG2_MAILER_HOST", "127.0.0.1")
	setEnv("CONFIG2_MAILER_USERNAME", "foo")
	setEnv("CONFIG2_MAILER_PASSWORD", "bar")
	setEnv("CONFIG2_MAILER_ENABLED", "true")
	setEnv("CONFIG2_MAP_KEY1", "value1")
	setEnv("CONFIG2_MAP_KEY2", "value2")
	setEnv("CONFIG2_MAP_KEY3", "value3")

	env, err := NewEnvTree("^CONFIG2", "_")

	if err != nil {
		return
	}

	s := CONFIG2{}

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
	// {"DB":{"USERNAME":"foo","PASSWORD":"bar","HOST":"localhost","NAME":"my-db","PORT":3306,"URL":"jdbc:mysql://localhost:3306/my-db?user=foo\u0026password=bar","USAGELIMIT":95.6},"MAILER":{"HOST":"127.0.0.1","USERNAME":"foo","PASSWORD":"bar","ENABLED":true},"MAP":{"KEY1":"value1","KEY2":"value2","KEY3":"value3"}}
}
```
