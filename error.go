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

// VariableDuplicatedError is triggered when environment variable cannot be found
type VariableDuplicatedError struct {
	Variable string
}

// Error dump error
func (e VariableDuplicatedError) Error() string {
	return fmt.Sprintf(`Variable "%s" was already defined before`, e.Variable)
}
