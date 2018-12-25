package logger

import "fmt"

//log level
const (
	Level_All = iota
	Level_Debug
	Level_Info
	Level_Warn
	Level_Error
	Level_Fatal
)

//the output type of logs
type OutputType byte

const (
	OutputType_all     OutputType = iota //both output to file and console
	OutputType_console                   //only output to console
	OutputType_file                      //only output to file
)

type SliceType byte

const (
	SliceType_Size SliceType = iota //slice file by file size
	SliceType_Date                  //slice file by date
)

func NewLogger(config string) (*Logger, error) {
	obj := &Logger{}
	err := obj.Init(config)
	return obj, err
}

type Logger struct {
	handler *logHandler
}

func (lg *Logger) Init(config string) error {
	if lg.handler == nil {
		lg.handler = new(logHandler)
		err := lg.handler.init(config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (lg *Logger) Log(level int, v ...interface{}) {
	data := fmt.Sprint(v...)
	lg.handler.inputData(level, data)
}

func (lg *Logger) LogF(level int, format string, a ...interface{}) {
	data := fmt.Sprintf(format, a...)
	lg.handler.inputData(level, data)
}

func (lg *Logger) Debug(v ...interface{}) {
	data := fmt.Sprint(v...)
	lg.handler.inputData(Level_Debug, data)
}

func (lg *Logger) DebugF(format string, a ...interface{}) {
	data := fmt.Sprintf(format, a...)
	lg.handler.inputData(Level_Debug, data)
}

func (lg *Logger) Info(v ...interface{}) {
	data := fmt.Sprint(v...)
	lg.handler.inputData(Level_Info, data)
}

func (lg *Logger) InfoF(format string, a ...interface{}) {
	data := fmt.Sprintf(format, a...)
	lg.handler.inputData(Level_Info, data)
}

func (lg *Logger) Warn(v ...interface{}) {
	data := fmt.Sprint(v...)
	lg.handler.inputData(Level_Warn, data)
}

func (lg *Logger) WarnF(format string, a ...interface{}) {
	data := fmt.Sprintf(format, a...)
	lg.handler.inputData(Level_Warn, data)
}

func (lg *Logger) Error(v ...interface{}) {
	data := fmt.Sprint(v...)
	lg.handler.inputData(Level_Error, data)
}

func (lg *Logger) ErrorF(format string, a ...interface{}) {
	data := fmt.Sprintf(format, a...)
	lg.handler.inputData(Level_Error, data)
}

func (lg *Logger) Fatal(v ...interface{}) {
	data := fmt.Sprint(v...)
	lg.handler.inputData(Level_Fatal, data)
}

func (lg *Logger) FatalF(format string, a ...interface{}) {
	data := fmt.Sprintf(format, a...)
	lg.handler.inputData(Level_Fatal, data)
}
