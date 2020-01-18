package flags

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func shift(ss []string) (string, []string) {
	return ss[0], ss[1:]
}

// Args creates a pair of empty positional and optional argument definitions.
func Args() (*Positional, *Optional) {
	return newPositional(), newOptional()
}

// Command represents a executable command.
type Command func(*Context) error

// CommandDescription carries a command and its description.
type CommandDescription struct {
	Desc string
	Cmd  Command
}

// Program represents a list of named commands.
type Program struct {
	Map map[string]CommandDescription
}

// NewProgram creates a new Program.
func NewProgram() *Program {
	return &Program{make(map[string]CommandDescription)}
}

// Add a Command with the given name and description.
func (prog *Program) Add(name, desc string, cmd Command) {
	prog.Map[name] = CommandDescription{desc, cmd}
}

// ListCommands lists the commands registered to the given program.
func (prog Program) ListCommands() string {
	names := make([]string, len(prog.Map))
	i := 0
	for name := range prog.Map {
		names[i] = name
		i++
	}
	sort.Strings(names)
	builder := strings.Builder{}
	builder.WriteString("available commands:")
	for _, name := range names {
		cmd := prog.Map[name]
		builder.WriteString("\n" + formatHelp(name, cmd.Desc))
	}
	return builder.String()
}

// Compile the subcommands into a single command.
func (prog Program) Compile() Command {
	return func(ctx *Context) error {
		if len(ctx.Args) == 0 {
			return fmt.Errorf("%s expected a command.\n\n%s", ctx.Name, prog.ListCommands())
		}
		head, tail := shift(ctx.Args)
		if strings.HasPrefix(head, "-h") || head == "--help" {
			return fmt.Errorf("%s: %s\n\n%s", ctx.Name, ctx.Desc, prog.ListCommands())
		}
		v, ok := prog.Map[head]
		if !ok {
			return fmt.Errorf("unknown command name `%s`", head)
		}
		name := fmt.Sprintf("%s %s", ctx.Name, head)
		return v.Cmd(&Context{name, v.Desc, tail})
	}
}

// Main is the main program.
var Main = NewProgram()

// Add a command to the main program.
func Add(name, desc string, cmd Command) { Main.Add(name, desc, cmd) }

// Compile the main program.
func Compile() Command { return Main.Compile() }

// Run the given command using os.Args.
func Run(name, desc string, cmd Command) int {
	ctx := &Context{name, desc, os.Args[1:]}
	if err := cmd(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
