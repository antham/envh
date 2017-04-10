package envh

import (
	"reflect"
)

type entry struct {
	typ   reflect.Type
	value reflect.Value
	chain []string
}

func populateInt(forceDefinition bool, tree *EnvTree, val reflect.Value, keyChain []string) error {
	v, err := tree.FindInt(keyChain...)

	if forceDefinition && err != nil {
		return err
	}

	if _, ok := err.(WrongTypeError); ok {
		return err
	}

	val.SetInt(int64(v))

	return nil
}

func populateFloat(forceDefinition bool, tree *EnvTree, val reflect.Value, keyChain []string) error {
	v, err := tree.FindFloat(keyChain...)

	if forceDefinition && err != nil {
		return err
	}

	if _, ok := err.(WrongTypeError); ok {
		return err
	}

	val.SetFloat(float64(v))

	return nil
}

func populateString(forceDefinition bool, tree *EnvTree, val reflect.Value, keyChain []string) error {
	v, err := tree.FindString(keyChain...)

	if forceDefinition && err != nil {
		return err
	}

	val.SetString(v)

	return nil
}

func populateBool(forceDefinition bool, tree *EnvTree, val reflect.Value, keyChain []string) error {
	v, err := tree.FindBool(keyChain...)

	if forceDefinition && err != nil {
		return err
	}

	if _, ok := err.(WrongTypeError); ok {
		return err
	}

	val.SetBool(v)

	return nil
}

func populateEntry(entries *[]entry, tree *EnvTree, forceDefinition bool) error {
	var err error

	typ := (*entries)[0].typ
	value := (*entries)[0].value
	chain := (*entries)[0].chain

	(*entries) = append([]entry{}, (*entries)[1:]...)

	for i := 0; i < typ.NumField(); i++ {
		val := value.Field(i)

		switch val.Type().Kind() {
		case reflect.Struct:
			*entries = append(*entries, entry{val.Type(), val, append(chain, typ.Field(i).Name)})
		case reflect.Int:
			err = populateInt(forceDefinition, tree, val, append(chain, typ.Field(i).Name))
		case reflect.Float32:
			err = populateFloat(forceDefinition, tree, val, append(chain, typ.Field(i).Name))
		case reflect.String:
			err = populateString(forceDefinition, tree, val, append(chain, typ.Field(i).Name))
		case reflect.Bool:
			err = populateBool(forceDefinition, tree, val, append(chain, typ.Field(i).Name))
		default:
			err = TypeUnsupported{val.Type().Kind().String(), "int32, float32, string, boolean or struct"}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func isPointerToStruct(data interface{}) bool {
	return !(reflect.TypeOf(data).Kind() != reflect.Ptr || reflect.TypeOf(data).Elem().Kind() != reflect.Struct)
}

func populateStructFromEnvTree(data interface{}, tree *EnvTree, forceDefinition bool) error {
	if !isPointerToStruct(data) {
		return TypeUnsupported{reflect.TypeOf(data).Kind().String(), "pointer to struct"}
	}

	entries := []entry{{reflect.TypeOf(data).Elem(), reflect.ValueOf(data).Elem(), []string{reflect.TypeOf(data).Elem().Name()}}}

	for {
		err := populateEntry(&entries, tree, forceDefinition)

		if err != nil {
			return err
		}

		if len(entries) == 0 {
			return nil
		}
	}
}
