package flags

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type ArgumentType int

const (
	LongType ArgumentType = iota
	ShortType
	ValueType
)

func TypeOf(s string) ArgumentType {
	if strings.HasPrefix(s, "--") && s != "--" {
		return LongType
	}
	if strings.HasPrefix(s, "-") && s != "-" {
		return ShortType
	}
	return ValueType
}

type Parser struct {
	Pos *Positional
	Opt *Optional
}

func NewParser(pos *Positional, opt *Optional) Parser {
	return Parser{pos, opt}
}

func (parser Parser) handleValue(name string, args []string) ([]string, error) {
	pos, opt := parser.Pos, parser.Opt
	head := ""

	switch v := opt.Args[name].Value.(type) {
	// Do not accept value arguments behind boolean flags.
	case *BoolValue:
		*v = BoolValue(true)

	case SliceValue:
		n := 0
		for _, arg := range args {
			if TypeOf(arg) == ValueType {
				n++
			}
		}

		for TypeOf(args[0]) == ValueType && n > pos.Len() {
			head, args = shift(args)
			v.Set(head)
			n--
		}

	default:
		head, args = shift(args)
		if TypeOf(head) != ValueType {
			return nil, fmt.Errorf("value not given for flag `--%s`", name)
		}
		v.Set(head)
	}

	return args, nil
}

var errHelp = errors.New("help")

func (parser Parser) Parse(args []string) error {
	pos, opt := parser.Pos, parser.Opt
	optmap := make(map[string]string)
	head := ""
	extra := []string{}

	for len(args) > 0 {
		head, args = shift(args)

		switch TypeOf(head) {

		// Process long flag name.
		case LongType:
			long := head[2:]

			if long == "help" {
				return errHelp
			}

			switch i := strings.IndexByte(head, '='); i {
			case -1:
				if !opt.Args.Has(long) {
					return fmt.Errorf("unknown flag `--%s`", long)
				}
				var err error
				args, err = parser.handleValue(long, args)
				if err != nil {
					return fmt.Errorf("in flag `--%s`: %v", long, err)
				}

			// Flag has form `--long=value`.
			default:
				name, value := long[:i], long[i+1:]
				if !opt.Args.Has(name) {
					return fmt.Errorf("unknown flag `--%s`", name)
				}
				optmap[name] = value
			}

		// Process short flag name.
		case ShortType:
			rr := []rune(head[1:])
			var r rune

			for len(rr) > 0 {
				r, rr = rr[0], rr[1:]

				if r == 'h' {
					return errHelp
				}

				name, ok := opt.Alias[r]
				if !ok {
					return fmt.Errorf("unknown shorthand `%c`", r)
				}

				switch len(rr) {
				// The last shorthand flag can be a non-boolean value
				case 0:
					var err error
					args, err = parser.handleValue(name, args)
					if err != nil {
						return fmt.Errorf("in flag `--%s`: %v", name, err)
					}

				default:
					switch v := opt.Args[name].Value.(type) {
					case *BoolValue:
						*v = BoolValue(true)
					default:
						return fmt.Errorf("flag `%s for shorthand `%c` is not boolean", name, r)
					}
				}
			}

		// The argument is not associated to a flag.
		default:
			extra = append(extra, head)
		}
	}

	for i, name := range pos.Order {
		if len(extra) == 0 {
			missing := strings.Join(pos.Order[i:], "`, `")
			return fmt.Errorf("missing positional argument(s): `%s`", missing)
		}
		head, extra = shift(extra)
		pos.Args[name].Value.Set(head)
	}

	if pos.In != nil && len(extra) > 0 {
		head, extra = shift(extra)
		info, err := os.Stat(head)
		if err == nil && !info.IsDir() {
			if err := pos.In.Value.Set(head); err != nil {
				return fmt.Errorf("in input file: %v", err)
			}
		}
	}

	if pos.Out != nil && len(extra) > 0 {
		head, extra = shift(extra)
		if err := pos.Out.Value.Set(head); err != nil {
			return fmt.Errorf("in output file: %v", err)
		}
	}

	if len(extra) > 0 {
		extras := strings.Join(extra, "`, `")
		return fmt.Errorf("extraneous arguments: `%s`", extras)
	}

	return nil
}
