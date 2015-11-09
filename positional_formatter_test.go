package logrus

import (
	_ "bytes"
	_ "errors"
	"fmt"
	"os"
	"runtime"
	"path/filepath"
	"testing"

)

func TestFunctionMap(t *testing.T) {
	functionsMap := make(map[string]func()string);
	functionsMap["linenum"] = func()string {
		_, path, line, _ :=  runtime.Caller(3)
		filename := filepath.Base(path)
		return fmt.Sprintf("%s:%d",filename,line)
	}

	f, _ := os.OpenFile("/tmp/out", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0744)
	defer f.Close()
	pf := &PositionalFormatter{

		Functions: functionsMap,
		Fields: []string{"level","linenum", "msg"},
	}
	/*log := &Logger{
		Out: f,
		Formatter: pf,
		Level: InfoLevel,
	}*/
	bytes, _ := pf.Format(&Entry{Level: InfoLevel, Message: "Hey"})
	t.Logf("%s", string(bytes) )


}

func TestEscapedString(t *testing.T) {
	pf := &PositionalFormatter{
		Fields: []string{"level","`LOG>", "msg"},
	}
	bytes, _ := pf.Format(&Entry{Level: InfoLevel, Message: "Hey"})
	expected := "info LOG> Hey"
	if string(bytes) != expected {
		t.Errorf("Output did not match. Expected: %s, Actual: %s", expected, string(bytes) )
	}
}
