package logrus

import (
	_ "bytes"
	_ "fmt"
	_ "runtime"
	_ "sort"
	_ "strings"
	_ "time"
)


// TODO
type PositionalFormatter struct {
	Functions map[string]func()
	Fields  []string
	TimestampFormat string
}

func (f *PositionalFormatter) Format(entry *Entry) ([]byte, error) {
	var keys []string = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys,k)
	}
}
