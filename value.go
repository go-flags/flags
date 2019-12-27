package flags

// Value represents a command line argument value.
type Value interface {
	Set(value string) error
	String() string
}

// SliceValue represents a variable length command line argument value.
type SliceValue interface {
	Value
	Len() int
}

// Argument represents a value-usages pair.
type Argument struct {
	Value Value
	Usage string
}

// Arguments is a map of names and arguments.
type Arguments map[string]Argument

// Has tests if the given name is registered as an argument.
func (args *Arguments) Has(name string) bool {
	_, ok := (*args)[name]
	return ok
}
