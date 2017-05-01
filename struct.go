package envh

import (
	"reflect"
)

// StructWalker must be implemented, when using PopulateStruct* functions,
// to be able to set a value for a custom field with an unsupported field (a map for instance),
// to add transformation before setting a field or for custom validation purpose.
// Walk function is called when struct is populated for every struct field a matching is made with
// an EnvTree node. Two parameters are given : tree represents whole parsed tree and keyChain is path leading to the node in tree.
// Returning true as first parameter will bypass walking process and false not, so it's
// possible to completely control how some part of a structure are defined and it's possible as well
// only to add some checking and let regular process do its job.
type StructWalker interface {
	Walk(tree *EnvTree, keyChain []string) (bypassWalkingProcess bool, err error)
}

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

func populateRegularType(entries *[]entry, tree *EnvTree, val reflect.Value, valKeyChain []string, forceDefinition bool) error {
	switch val.Type().Kind() {
	case reflect.Struct:
		*entries = append(*entries, entry{val.Type(), val, valKeyChain})

		return nil
	case reflect.Int:
		return populateInt(forceDefinition, tree, val, valKeyChain)
	case reflect.Float32:
		return populateFloat(forceDefinition, tree, val, valKeyChain)
	case reflect.String:
		return populateString(forceDefinition, tree, val, valKeyChain)
	case reflect.Bool:
		return populateBool(forceDefinition, tree, val, valKeyChain)
	default:
		return TypeUnsupported{val.Type().Kind().String(), "int32, float32, string, boolean or struct"}
	}
}

func callStructMethodWalk(origStruct interface{}, tree *EnvTree, keyChain []string) (bool, error) {
	if walker, ok := origStruct.(StructWalker); ok {
		return walker.Walk(tree, keyChain)
	}

	return false, nil
}

func populateStruct(entries *[]entry, origStruct interface{}, tree *EnvTree, forceDefinition bool) error {
	var err error
	var ok bool
	var val reflect.Value
	var valKeyChain []string

	typ := (*entries)[0].typ
	value := (*entries)[0].value
	chain := (*entries)[0].chain

	(*entries) = append([]entry{}, (*entries)[1:]...)

	for i := 0; i < typ.NumField(); i++ {
		val = value.Field(i)
		valKeyChain = append([]string{}, append(chain, typ.Field(i).Name)...)

		ok, err = callStructMethodWalk(origStruct, tree, valKeyChain)

		if err != nil {
			return err
		}

		if ok {
			continue
		}

		if err = populateRegularType(entries, tree, val, valKeyChain, forceDefinition); err != nil {
			return err
		}
	}

	return nil
}

func isPointerToStruct(data interface{}) bool {
	return !(reflect.TypeOf(data).Kind() != reflect.Ptr || reflect.TypeOf(data).Elem().Kind() != reflect.Struct)
}

func populateStructFromEnvTree(origStruct interface{}, tree *EnvTree, forceDefinition bool) error {
	if !isPointerToStruct(origStruct) {
		return TypeUnsupported{reflect.TypeOf(origStruct).Kind().String(), "pointer to struct"}
	}

	entries := []entry{{reflect.TypeOf(origStruct).Elem(), reflect.ValueOf(origStruct).Elem(), []string{reflect.TypeOf(origStruct).Elem().Name()}}}

	for {
		err := populateStruct(&entries, origStruct, tree, forceDefinition)

		if err != nil {
			return err
		}

		if len(entries) == 0 {
			return nil
		}
	}
}
