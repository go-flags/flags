package flags

import (
	"fmt"
	"os"
)

// Positional represents the positional command line arguments.
type Positional struct {
	Order []string
	Args  Arguments
	In    *Argument
	Out   *Argument
}

func newPositional() *Positional {
	return &Positional{[]string{}, Arguments{}, nil, nil}
}

// Len returns the number of positional arguments.
func (pos Positional) Len() int { return len(pos.Args) }

// Regsiter the name with the given value and usage.
func (pos *Positional) Register(name string, value Value, usage string) {
	if pos.Args.Has(name) {
		panic(fmt.Errorf("positional argument with name `%s`already exists", name))
	}
	pos.Order = append(pos.Order, name)
	pos.Args[name] = Argument{value, usage}
}

// Bool adds a string value to the positional argument list.
func (pos *Positional) Bool(name, usage string) *bool {
	value := NewBoolValue(false)
	pos.Register(name, value, usage)
	return (*bool)(value)
}

// Int adds a string value to the positional argument list.
func (pos *Positional) Int(name, usage string) *int {
	value := NewIntValue(0)
	pos.Register(name, value, usage)
	return (*int)(value)
}

// String adds a string value to the positional argument list.
func (pos *Positional) String(name, usage string) *string {
	value := NewStringValue("")
	pos.Register(name, value, usage)
	return (*string)(value)
}

// Open adds a file for reading to the positional argument list.
func (pos *Positional) Open(name, usage string) *os.File {
	value := NewOpenValue(nil)
	pos.Register(name, value, usage)
	return (*os.File)(value)
}

// Create adds a file for writing to the positional argument list.
func (pos *Positional) Create(name, usage string) *os.File {
	value := NewCreateValue(nil)
	pos.Register(name, value, usage)
	return (*os.File)(value)
}

// Input adds a file which when omitted will read from os.Stdin.
func (pos *Positional) Input(usage string) *os.File {
	value := NewOpenValue(os.Stdin)
	pos.In = &Argument{value, usage}
	return (*os.File)(value)
}

// Output adds a file which when omitted will read from os.Stdout.
func (pos *Positional) Output(usage string) *os.File {
	value := NewCreateValue(os.Stdout)
	pos.Out = &Argument{value, usage}
	return (*os.File)(value)
}
