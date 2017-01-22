package envh

import (
	"os"

	"github.com/Sirupsen/logrus"
)

func setTestingEnvs() {
	datas := map[string]string{
		"TEST1": "test1",
		"TEST2": "=test2=",
	}

	for k, v := range datas {
		err := os.Setenv(k, v)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func setTestingEnvsForTree() {
	datas := map[string]string{
		"ENVH_TEST1_TEST2_TEST3": "test1",
		"ENVH_TEST1_TEST2_TEST4": "test2",
		"ENVH_TEST1_TEST5_TEST6": "test3",
		"ENVH_TEST1_TEST7_TEST2": "test4",
		"ENVH_TEST1":             "test5",
	}

	for k, v := range datas {
		err := os.Setenv(k, v)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func setEnv(key string, value string) {
	err := os.Setenv(key, value)

	if err != nil {
		logrus.Fatal(err)
	}
}
