package envh

import (
	"os"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
)

var envs []string

func TestMain(m *testing.M) {
	saveExistingEnvs()
	code := m.Run()
	os.Exit(code)
}

func saveExistingEnvs() {
	envs = os.Environ()
}

func setEnv(key string, value string) {
	err := os.Setenv(key, value)

	if err != nil {
		logrus.Fatal(err)
	}
}

func restoreEnvs() {
	os.Clearenv()

	if len(envs) != 0 {
		for _, envCouple := range envs {
			parseEnv := strings.Split(envCouple, "=")

			setEnv(parseEnv[0], strings.Join(parseEnv[1:], "="))
		}
	}
}
