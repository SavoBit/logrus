package logrus

import "fmt"

type CustomJSONFormatter struct {
	Functions map[string]func() string
	Fields    []string
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
}

func (f *CustomJSONFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(Fields, len(entry.Data)+3+len(f.Fields))
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	data["time"] = entry.Time.Format(timestampFormat)
	data["msg"] = entry.Message
	data["level"] = entry.Level.String()
	if len(f.Fields) > 0 {
		for _, key := range f.Fields {
			switch {
			case !(key == "time" || key == "msg" || key == "level"):
				if function, ok := f.Functions[key]; ok {
					val := function()
					data[key] = val
				}
			}
		}
	}

	serialized, err := jsonMarshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
