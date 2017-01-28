package envh

import (
	"strconv"
)

func getString(fun func() (string, bool)) (string, error) {
	if v, ok := fun(); ok {
		return v, nil
	}

	return "", VariableNotFoundError{}
}

func getInt(fun func() (string, bool)) (int, error) {
	v, ok := fun()

	if !ok {
		return 0, VariableNotFoundError{}
	}

	i, err := strconv.Atoi(v)

	if err != nil {
		return 0, WrongTypeError{v, "int"}
	}

	return i, nil
}

func getFloat(fun func() (string, bool)) (float32, error) {
	v, ok := fun()

	if !ok {
		return 0, VariableNotFoundError{}
	}

	f, err := strconv.ParseFloat(v, 32)

	if err != nil {
		return 0, WrongTypeError{v, "float"}
	}

	return float32(f), nil
}

func getBool(fun func() (string, bool)) (bool, error) {
	v, ok := fun()

	if !ok {
		return false, VariableNotFoundError{}
	}

	b, err := strconv.ParseBool(v)

	if err != nil {
		return false, WrongTypeError{v, "bool"}
	}

	return b, nil
}
