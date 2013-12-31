package shell

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
)

func NewFlagSet(name string) *flag.FlagSet {
	return flag.NewFlagSet(name, flag.ContinueOnError)
}

func CommandWord(command string) string {
	if i := strings.Index(command, " "); i != -1 {
		return command[:i]
	}
	return command
}

func Parse(fs *flag.FlagSet, command string) (string, error) {
	r := make([]string, 0)
	var st int = 0 // 0-normal, 1-"' 2-$
	var pos int = 0
	var sp rune = 0

	for i, ch := range []rune(command) {

		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			switch st {
			case 0, 2:
				st = 0
				if pos < i {
					r = append(r, command[pos:i])
				}
				pos = i + 1
			case 1:
			}
			continue
		}

		if ch == '"' || ch == '\'' {
			switch st {
			case 0, 2:
				if pos < i {
					r = append(r, command[pos:i])
				}
				pos = i + 1
				st = 1
				sp = ch
			case 1:
				if sp == ch {
					if pos < i {
						r = append(r, command[pos:i])
						pos = i + 1
					}
					st = 0
				}
			}
			continue
		}

		if ch == '$' {
			end := false
			switch st {
			case 0:
				if i-pos == 0 {
					st = 2
				}
			case 1:
			case 2:
				pos = i + 1
				end = true
			}
			if end {
				break
			}
		}
	}
	if pos < len(command) {
		r = append(r, command[pos:])
	}
	if len(r) > 1 {
		return r[0], fs.Parse(r[1:])
	} else {
		return r[0], fs.Parse([]string{})
	}

}

func GetHelp(cmd string, args string, fs *flag.FlagSet) string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("usage: ")
	buf.WriteString(cmd)
	buf.WriteString("[options] ")
	buf.WriteString(args)

	fs.VisitAll(func(f *flag.Flag) {
		buf.WriteString("\n ")
		buf.WriteString(f.Name)
		if f.DefValue != "" {
			buf.WriteString(" (")
			buf.WriteString(f.DefValue)
			buf.WriteString(") ")
		}
		if f.Usage != "" {
			buf.WriteString(" ")
			buf.WriteString(f.Usage)
		}
	})
	return buf.String()
}

func DoParse(s *Session, command string, fs *flag.FlagSet, name string, args string, minargs int, maxargs int) bool {
	var help bool
	fs.BoolVar(&help, "h", false, "show help")
	_, err := Parse(fs, command)
	if err != nil {
		s.Writeln(fmt.Sprintf("ERROR: %s", err))
		return true
	}
	if help {
		s.Writeln(GetHelp(name, args, fs))
		return true
	}
	if minargs >= 0 && fs.NArg() < minargs {
		s.Writeln(fmt.Sprintf("ERROR: args < %d", minargs))
		s.Writeln(GetHelp(name, args, fs))
		return true
	}
	if maxargs >= 0 && fs.NArg() > maxargs {
		s.Writeln(fmt.Sprintf("ERROR: args > %d", maxargs))
		s.Writeln(GetHelp(name, args, fs))
		return true
	}
	return false
}
