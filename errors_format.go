package errors

import (
	"fmt"
	"io"
	"strings"
)

func (e *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(e.msg) > 0 {
				io.WriteString(s, e.msg)
				if e.msg[len(e.msg)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "stack trace (most recent call first):\n")
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}

func (e *msgWithStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(e.msg) > 0 {
				io.WriteString(s, e.msg)
				if e.msg[len(e.msg)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "stack trace (most recent call first):\n")
			stack(e.StackTrace()).Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}

func (e *msgWithoutStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(e.msg) > 0 {
				io.WriteString(s, e.msg)
				if e.msg[len(e.msg)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "stack trace (most recent call first):\n")
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}

func (e *msgJoined) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(e.msg) > 0 {
				io.WriteString(s, e.msg)
				if e.msg[len(e.msg)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "multiple errors at (most recent call first):\n")
			e.stack.Format(s, verb)
			for _, err := range e.errs {
				unwrap := fmt.Sprintf("%+v", err)
				unwrap = strings.TrimRight(unwrap, "\n")
				lines := strings.Split(unwrap, "\n")
				for i := range lines {
					lines[i] = "\t" + lines[i]
				}
				io.WriteString(s, "\n")
				io.WriteString(s, strings.Join(lines, "\n"))
				io.WriteString(s, "\n")
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}

func (e *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", e.err)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *withoutStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			unwrap := fmt.Sprintf("%+v", e.err)
			if len(unwrap) > 0 {
				io.WriteString(s, unwrap)
				if unwrap[len(unwrap)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "stack trace (most recent call first):\n")
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *cause) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(e.msg) > 0 {
				io.WriteString(s, e.msg)
				if e.msg[len(e.msg)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "stack trace (most recent call first):\n")
			e.stack.Format(s, verb)
			unwrap := fmt.Sprintf("%+v", e.err)
			if len(unwrap) > 0 {
				io.WriteString(s, "\nthe above error was caused by the following error:\n\n")
				io.WriteString(s, unwrap)
				if unwrap[len(unwrap)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}

func (e *joined) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "multiple errors at (most recent call first):\n")
			e.stack.Format(s, verb)
			for _, err := range e.errs {
				unwrap := fmt.Sprintf("%+v", err)
				unwrap = strings.TrimRight(unwrap, "\n")
				lines := strings.Split(unwrap, "\n")
				for i := range lines {
					lines[i] = "\t" + lines[i]
				}
				io.WriteString(s, "\n")
				io.WriteString(s, strings.Join(lines, "\n"))
				io.WriteString(s, "\n")
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}
