package log

import (
	"bufio"
	"fmt"
	"io"
	"runtime"
	"time"
)

const (
	logfmtTimeFormat = "2006-01-02T15:04:05.000000-0700"
	timeFormat       = "2006-01-02T15:04:05-0700"
)

var eol string

func init() {
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	} else {
		eol = "\n"
	}
}

// Message contains all of the log message information.
// Note that *Message implements the error interface.
type Message struct {
	Timestamp  time.Time
	Level      Level
	Text       string
	Err        error
	Parameters []Parameter
	Context    []Parameter
	errorCode  string
	statusCode int
}

// Parameter contains additional information about
// a log message.
type Parameter struct {
	Name  string
	Value interface{}
}

func newMessage(level Level, text string) *Message {
	m := &Message{
		Timestamp: time.Now(),
		Level:     level,
		Text:      text,
	}
	return m
}

func (m *Message) applyOpt(opt Option) *Message {
	opt(m)
	return m
}

// applyOpts applies all of the option functions to the message.
func (m *Message) applyOpts(opts []Option) *Message {
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Error implements the error interface
func (m *Message) Error() string {
	return m.Text
}

// Code returns any error code string associated with the message.
func (m *Message) Code() string {
	return m.errorCode
}

// StatusCode returns any HTTP status code associated with the message.
// When loggging a message, it is possible to suggest an appropriate HTTP status
// code to return to the client.
func (m *Message) StatusCode() int {
	return m.statusCode
}

func writeValueString(w *bufio.Writer, value string) error {
	var err error
	needsQuotes := false
	needsEscape := false
	hasEscape := false
	for _, c := range value {
		if c == '"' || c == '=' {
			needsQuotes = true
			needsEscape = true
		} else if c <= ' ' {
			needsQuotes = true
		} else if c == '\\' {
			hasEscape = true
		}
	}

	if needsQuotes && hasEscape {
		needsEscape = true
	}

	if needsQuotes {
		_, err = w.WriteRune('"')
		if err != nil {
			return err
		}
		if needsEscape {
			// need to escape one or more chars in the value
			for _, c := range value {
				if c == '\\' || c == '"' {
					_, err = w.WriteRune('\\')
					if err != nil {
						return err
					}
				}
				_, err = w.WriteRune(c)
				if err != nil {
					return err
				}
			}
		} else {
			// no escapable chars in the value, so can just write the contents
			_, err := w.WriteString(value)
			if err != nil {
				return err
			}
		}
		_, err := w.WriteRune('"')
		if err != nil {
			return err
		}
	} else {
		// no quotes required, just write the string
		_, err := w.WriteString(value)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeProperty(w *bufio.Writer, key string, value interface{}) error {
	var err error
	w.WriteRune(' ')
	w.WriteString(key)
	w.WriteRune('=')
	switch v := value.(type) {
	case bool:
		_, err = fmt.Fprint(w, v)
		return err
	case byte:
		_, err = fmt.Fprint(w, v)
		return err
	case complex64:
		return writeValueString(w, fmt.Sprint(v))
	case complex128:
		return writeValueString(w, fmt.Sprint(v))
	case error:
		return writeValueString(w, v.Error())
	case float32:
		_, err = fmt.Fprint(w, v)
		return err
	case float64:
		_, err = fmt.Fprint(w, v)
		return err
	case int:
		_, err = fmt.Fprint(w, v)
		return err
	case int16:
		_, err = fmt.Fprint(w, v)
		return err
	case int32:
		_, err = fmt.Fprint(w, v)
		return err
	case int64:
		_, err = fmt.Fprint(w, v)
		return err
	case int8:
		_, err = fmt.Fprint(w, v)
		return err
	case string:
		return writeValueString(w, v)
	case uint:
		_, err = fmt.Fprint(w, v)
		return err
	case uint16:
		_, err = fmt.Fprint(w, v)
		return err
	case uint32:
		_, err = fmt.Fprint(w, v)
		return err
	case uint64:
		_, err = fmt.Fprint(w, v)
		return err
	case uintptr:
		_, err = fmt.Fprint(w, v)
		return err
	}

	if v, ok := value.(fmt.Stringer); ok {
		return writeValueString(w, v.String())
	}

	return writeValueString(w, fmt.Sprint(value))
}

func (m *Message) Logfmt(w io.Writer) {
	writer := bufio.NewWriterSize(w, 1024)
	writer.WriteString(m.Timestamp.Format(logfmtTimeFormat))
	writer.WriteRune(' ')
	writer.WriteString(m.Level.String())
	writeProperty(writer, "msg", m.Text)
	if m.Err != nil {
		writeProperty(writer, "error", m.Err.Error())
	}
	for _, p := range m.Parameters {
		writeProperty(writer, p.Name, p.Value)
	}
	for _, p := range m.Context {
		writeProperty(writer, p.Name, p.Value)
	}
	if m.errorCode != "" {
		writeProperty(writer, "code", m.errorCode)
	}
	if m.statusCode != 0 {
		writeProperty(writer, "status", m.statusCode)
	}
	writer.WriteString(eol)
	writer.Flush()
}

// WriteTo prints the log Text Message to the writer w.
func (m *Message) WriteTo(w io.Writer) {
	m.Logfmt(w)
	return

	// append parameters and context to form one larger list
	parameters := make([]Parameter, 0, len(m.Parameters)+len(m.Context)+1)
	parameters = append(parameters, m.Parameters...)
	if m.Err != nil {
		parameters = append(parameters, Parameter{"error", m.Err.Error()})
	}
	parameters = append(parameters, m.Context...)

	// TODO: more sanitizing of parameter values, particularly
	// strings as they might be malicious client input
	switch len(parameters) {
	case 0, 1, 2, 3, 4:
		fmt.Fprintf(w, "%s %-5s %s", m.Timestamp.Format(timeFormat), m.Level, m.Text)
		for _, param := range parameters {
			switch v := param.Value.(type) {
			case string:
				fmt.Fprintf(w, ", %s=%q", param.Name, v)
			case *string:
				if v == nil {
					fmt.Fprintf(w, ", %s=nil", param.Name)
				} else {
					fmt.Fprintf(w, ", %s=%q", param.Name, *v)
				}
			default:
				fmt.Fprintf(w, ", %s=%+v", param.Name, param.Value)
			}
		}
		fmt.Fprintf(w, "\n")
	default:
		fmt.Fprintf(w, "%s %-5s %s\n", m.Timestamp.Format(timeFormat), m.Level, m.Text)
		for _, param := range parameters {
			switch v := param.Value.(type) {
			case string:
				fmt.Fprintf(w, "    %s=%q\n", param.Name, v)
			case *string:
				if v == nil {
					fmt.Fprintf(w, "    %s=nil\n", param.Name)
				} else {
					fmt.Fprintf(w, "    %s=%q\n", param.Name, *v)
				}
			default:
				fmt.Fprintf(w, "    %s=%+v\n", param.Name, param.Value)
			}
		}
	}
}
