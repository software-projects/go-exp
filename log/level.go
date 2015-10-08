package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Level indicates the level of a log message.
type Level int

const (
	LevelDebug   Level = iota // Debugging only
	LevelInfo                 // Informational
	LevelWarning              // Warning
	LevelError                // Error condition
)

// MinLevel is the minimum level that will be logged.
// The calling program can change this value at any time.
var (
	MinLevel = LevelInfo
)

var (
	errInvalidLevel = errors.New("invalid level")
)

// String implements the String interface.
func (s Level) String() string {
	switch s {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warn"
	case LevelError:
		return "error"
	}
	return fmt.Sprintf("unknown %d", s)
}

// MarshalJSON implements the json.Marshaler interface.
func (s Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Level) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return errInvalidLevel
	}
	switch strings.ToLower(str) {
	case "debug":
		*s = LevelDebug
	case "info":
		*s = LevelInfo
	case "warn", "warning":
		*s = LevelWarning
	case "error":
		*s = LevelError
	default:
		return errInvalidLevel
	}
	return nil
}
