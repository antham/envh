package envh

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type CONFIG3 struct {
	SERVER1 struct {
		IP   string
		PORT string
	}
	SERVER2 struct {
		IP   string
		PORT string
	}
}

func (c *CONFIG3) Walk(tree *EnvTree, keyChain []string) (bool, error) {
	if validator, ok := map[string]func(*EnvTree, []string) error{
		"CONFIG3_SERVER1_IP":   c.validateIP,
		"CONFIG3_SERVER2_IP":   c.validateIP,
		"CONFIG3_SERVER1_PORT": c.validatePort,
		"CONFIG3_SERVER2_PORT": c.validatePort,
	}[strings.Join(keyChain, "_")]; ok {
		return false, validator(tree, keyChain)
	}

	return false, nil
}

func (c *CONFIG3) validateIP(tree *EnvTree, keyChain []string) error {
	ipStr, err := tree.FindString(keyChain...)

	if err != nil {
		return err
	}

	if ip := net.ParseIP(ipStr); ip == nil {
		return fmt.Errorf(`"%s" is not a valid IP change "%s"`, ipStr, strings.Join(keyChain, "_"))
	}

	return nil
}

func (c *CONFIG3) validatePort(tree *EnvTree, keyChain []string) error {
	port, err := tree.FindInt(keyChain...)

	if err != nil {
		return err
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf(`"%d" is not a valid port, must be comprised between 1 and 65535 "%s"`, port, strings.Join(keyChain, "_"))
	}

	return nil
}

func ExampleStructWalker_customValidation() {
	os.Clearenv()
	setEnv("CONFIG3_SERVER1_IP", "127.0.0.1")
	setEnv("CONFIG3_SERVER1_PORT", "3000")
	setEnv("CONFIG3_SERVER2_IP", "localhost")
	setEnv("CONFIG3_SERVER2_PORT", "4000")

	env, err := NewEnvTree("^CONFIG3", "_")

	if err != nil {
		return
	}

	s := CONFIG3{}

	err = env.PopulateStruct(&s)

	fmt.Println(err)
	// Output: "localhost" is not a valid IP change "CONFIG3_SERVER2_IP"
}
