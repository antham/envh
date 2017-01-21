package envh

import (
	"fmt"
)

// ErrNotFound is triggered when environment variable cannot be found
var ErrNotFound = fmt.Errorf("Variable not found")

// ErrWrongType is triggered when we try to convert variable to a wrong type
var ErrWrongType = fmt.Errorf("Variable can't be converted")

// ErrDuplicated is triggered when a variable is already defined in a tree structure
var ErrDuplicated = fmt.Errorf("Variable was already defined before")
