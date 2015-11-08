package logrus

import (
	_ "bytes"
	_ "errors"
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestFunctionMap(t *testing.T) {
	functionsMap := make(map[string]func()string);
	functionsMap["linenum"] = func()string {
			pc := make([]uintptr, 10)
			runtime.Callers(2, pc)
			f := runtime.FuncForPC(pc[0])
			file, line := f.FileLine(pc[0])
		return fmt.Sprintf("%s:%d",file,line)
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
