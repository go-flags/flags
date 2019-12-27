package flags

import (
	"fmt"
	"os"
)

var shortNames = []rune("#%123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz")

type optionalName struct {
	Short rune
	Long  string
}

type byShort []optionalName

func (names byShort) Len() int { return len(names) }

func (names byShort) Less(i, j int) bool {
	a, b := names[i], names[j]
	switch {
	case a.Short != 0 && b.Short != 0:
		return a.Short < b.Short
	case a.Short != 0:
		return a.Short < []rune(b.Long)[0]
	case b.Short != 0:
		return []rune(a.Long)[0] < b.Short
	default:
		return a.Long < b.Long
	}
}

func (names byShort) Swap(i, j int) {
	names[i], names[j] = names[j], names[i]
}

// Optional represents the optional command line arguments.
type Optional struct {
	Args  Arguments
	Alias map[rune]string
}

func newOptional() *Optional {
	return &Optional{Arguments{}, make(map[rune]string)}
}

// Optional represents the optional command line arguments.
func (opt *Optional) Register(short rune, long string, value Value, usage string) {
	if opt.Args.Has(long) {
		panic(fmt.Errorf("optional argument with long name `%s` already exists", long))
	}
	if _, ok := opt.Alias[short]; ok {
		panic(fmt.Errorf("optional argument with short name `%c` already exists for long name `%s`", short, long))
	}
	if short != 0 {
		opt.Alias[short] = long
	}
	opt.Args[long] = Argument{value, usage}
}

// Switch adds a command line switch to the optional argument list.
func (opt *Optional) Switch(short rune, long string, usage string) *bool {
	value := NewBoolValue(false)
	opt.Register(short, long, value, usage)
	return (*bool)(value)
}

// Int adds an integer flag to the optional argument list.
func (opt *Optional) Int(short rune, long string, init int, usage string) *int {
	value := NewIntValue(init)
	opt.Register(short, long, value, usage)
	return (*int)(value)
}

// Float adds an float flag to the optional argument list.
func (opt *Optional) Float(short rune, long string, init float64, usage string) *float64 {
	value := NewFloatValue(init)
	opt.Register(short, long, value, usage)
	return (*float64)(value)
}

// String adds a string flag to the optional argument list.
func (opt *Optional) String(short rune, long, init, usage string) *string {
	value := NewStringValue(init)
	opt.Register(short, long, value, usage)
	return (*string)(value)
}

// StringSlice adds a string slice flag to the optional argument list.
func (opt *Optional) StringSlice(short rune, long string, init []string, usage string) *[]string {
	value := NewStringSliceValue(init)
	opt.Register(short, long, value, usage)
	return (*[]string)(value)
}

// OpenSlice adds a string slice flag to the optional argument list.
func (opt *Optional) OpenSlice(short rune, long string, init []*os.File, usage string) *[]*os.File {
	value := NewOpenSliceValue(init)
	opt.Register(short, long, value, usage)
	return (*[]*os.File)(value)
}
