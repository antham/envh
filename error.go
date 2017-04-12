package envh

import (
	"fmt"
	"strings"
)

// VariableNotFoundError is triggered when environment variable cannot be found
type VariableNotFoundError struct {
}

// Error dump error
func (e VariableNotFoundError) Error() string {
	return "Variable not found"
}

// NodeNotFoundError is triggered when tree node cannot be found
type NodeNotFoundError struct {
	KeyChain []string
}

// Error dump error
func (e NodeNotFoundError) Error() string {
	return fmt.Sprintf(`No node found at path "%s"`, strings.Join(e.KeyChain, " -> "))
}

// WrongTypeError is triggered when we try to convert variable to a wrong type
type WrongTypeError struct {
	Value interface{}
	Type  string
}

// Error dump error
func (e WrongTypeError) Error() string {
	return fmt.Sprintf(`Value "%s" can't be converted to type "%s"`, e.Value, e.Type)
}

// TypeUnsupported is triggered when a type isn't supported
type TypeUnsupported struct {
	ActualType   string
	RequiredType string
}

// Error dump error
func (e TypeUnsupported) Error() string {
	return fmt.Sprintf(`Type "%s" is not supported : you must provide "%s"`, e.ActualType, e.RequiredType)
}
