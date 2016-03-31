package logrus

import (
	"bytes"
	"fmt"
	"strings"
)

// PositionalFormatter is an attempt to support more "classical" style
// log formatters. The fields are written in order. In addition, a map
// of functions is supplied so that when the key is encountered, the
// corresponding function can be called and placed in the log
// this would be useful for supporting file name, line numbers, etc
type PositionalFormatter struct {
	Functions       map[string]func() string
	Fields          []string
	TimestampFormat string
	MultiLine       bool
}

// Format implements the formatter
func (f *PositionalFormatter) Format(entry *Entry) ([]byte, error) {
	var keys []string = make([]string, 0, len(entry.Data))
	for k := range entry.Data {

		keys = append(keys, k)
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	b := &bytes.Buffer{}

	for _, field := range f.Fields {

		switch {
		case field == "time":
			f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat))
		case field == "level":
			f.appendKeyValue(b, "level", entry.Level.String())
		case field == "msg":
			if f.MultiLine {
				b.WriteString(entry.Message)
			} else {
				f.appendKeyValue(b, "msg", entry.Message)
			}
		case field == "linenum":
			b.WriteString(entry.Data["linenum"].(string) + " ")
		case field == "package":
			f.appendKeyValue(b, "package", entry.Data["package"])
		case strings.HasPrefix(field, "`"):
			f.appendKeyValue(b, field, field)
		default:
			f.appendKeyValue(b, field, nil)
		}

	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func writeQuotedValue(value interface{}, b *bytes.Buffer) {
	switch value := value.(type) {
	case string:
		if needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	case error:
		errmsg := value.Error()
		if needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	default:
		fmt.Fprint(b, value)
	}
}

func (f *PositionalFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	// If the key has been overridden in functions, run that
	if function, ok := f.Functions[key]; ok {
		val := function()
		writeQuotedValue(val, b)
	} else {
		if strings.HasPrefix(key, "`") {
			b.WriteString(key[1:len(key)])
		} else {

			writeQuotedValue(value, b)
		}

	}
	b.WriteByte(' ')
}
