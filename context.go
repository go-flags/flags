package flags

import (
	"fmt"

	wrap "gopkg.in/ktnyt/wrap.v1"
)

// Context carries the name, description, and arguments given to a command.
type Context struct {
	Name string
	Desc string
	Args []string
}

// Parse the context arguments using the positional and optional argument
// definitions given.
func (ctx Context) Parse(pos *Positional, opt *Optional) error {
	parser := Parser{pos, opt}
	if err := parser.Parse(ctx.Args); err != nil {
		name := ctx.Name
		usage := wrap.Space(Usage(pos, opt), 72-len(name))
		if err == errHelp {
			return fmt.Errorf("usage: %s %s\n%s", ctx.Name, usage, Help(pos, opt))
		}
		return fmt.Errorf("%v\nusage: %s %s", err, ctx.Name, usage)
	}
	return nil
}

// Raise creates an error with the current context.
func (ctx Context) Raise(err error) error {
	if err != nil {
		return fmt.Errorf("%s: %v", ctx.Name, err)
	}
	return nil
}
