package logger

import "fmt"

const (
	Level_All = iota
	Level_Debug
	Level_Info
	Level_Warn
	Level_Error
	Level_Fatal
)

type SliceType byte

const (
	SliceType_Size SliceType = iota
	SliceType_Date
)

func NewLogger(config string) (*Logger, error) {
	obj := &Logger{}
	err := obj.init(config)
	return obj, err
}

type Logger struct {
	handler *logHandler
}

func (lg *Logger) init(config string) error {
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
