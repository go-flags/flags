package flags

import (
	"fmt"
	"sort"
	"strings"

	wrap "gopkg.in/ktnyt/wrap.v1"
)

func formatHelp(name, desc string) string {
	desc = wrap.Space(desc, 55)
	desc = strings.ReplaceAll(desc, "\n", "\n                        ")
	if len(name) < 22 {
		return "  " + name + strings.Repeat(" ", 22-len(name)) + desc
	}
	return "  " + name + "\n                        " + desc
}

// Usage creates a usage string for the given argument definitions.
func Usage(pos *Positional, opt *Optional) string {
	builder := strings.Builder{}
	builder.WriteString("[-h | --help]")
	if opt != nil && len(opt.Args) > 0 {
		builder.WriteString(" [<args>]")
	}
	if pos != nil {
		for _, name := range pos.Order {
			builder.WriteString(fmt.Sprintf(" <%s>", name))
		}
		if pos.In != nil {
			builder.WriteString(" [<infile>]")
		}
		if pos.Out != nil {
			builder.WriteString(" [<outfile>]")
		}
	}
	return builder.String()
}

// Help creaes a help string for the given argument definitions.
func Help(pos *Positional, opt *Optional) string {
	parts := []string{}
	if pos != nil {
		parts = append(parts, "\npositional arguments")
		for _, name := range pos.Order {
			usage := pos.Args[name].Usage
			name = fmt.Sprintf("<%s>", name)
			parts = append(parts, formatHelp(name, usage))
		}
		if pos.In != nil {
			usage := wrap.Space(pos.In.Usage, 55)
			parts = append(parts, formatHelp("[<infile>]", usage))
		}
		if pos.Out != nil {
			usage := wrap.Space(pos.Out.Usage, 55)
			parts = append(parts, formatHelp("[<outfile>]", usage))
		}
	}
	if opt != nil {
		parts = append(parts, "\noptional arguments:")
		names := []optionalName{}
		for long := range opt.Args {
			name := optionalName{0, long}
			for short := range opt.Alias {
				if opt.Alias[short] == long {
					name.Short = short
				}
			}
			names = append(names, name)
		}
		sort.Sort(byShort(names))
		for _, name := range names {
			long, short := name.Long, name.Short
			arg := opt.Args[long]
			usage := fmt.Sprintf("%s (value: %s)", arg.Usage, arg.Value)
			flag := ""
			switch arg.Value.(type) {
			case *BoolValue:
				flag = "--" + long
				if short != 0 {
					flag = fmt.Sprintf("-%c, %s", short, flag)
				}
			case SliceValue:
				flags := []string{}
				if short != 0 {
					flags = append(flags,
						fmt.Sprintf("-%[1]c <%[2]s> [<%[2]s> ...]", short, long),
						fmt.Sprintf("-%[1]c <%[2]s> [-%[1]c <%[2]s> ...]", short, long),
					)
				}
				flags = append(flags,
					fmt.Sprintf("--%[1]s <%[1]s> [<%[1]s> ...]", long),
					fmt.Sprintf("--%[1]s <%[1]s> [--%[1]s <%[1]s> ...]", long),
				)
				flag = strings.Join(flags, ",\n  ")
			default:
				flag = fmt.Sprintf("--%[1]s <%[1]s>", long)
				if short != 0 {
					flag = fmt.Sprintf("-%c <%s>, %s", short, long, flag)
				}
			}
			parts = append(parts, formatHelp(flag, usage))
		}
	}
	return strings.Join(parts, "\n")
}
