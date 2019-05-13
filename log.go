package dsky

type logModeType string
type LogAction string

const (
	logModeTypeInfo  logModeType = "info"
	logModeTypeDebug             = "debug"
	logModeTypeWarn              = "warn"
	logModeTypeError             = "error"

	LogActionDone LogAction = "done"
	LogActionWait           = "wait"
	LogActionFail           = "fail"
)

type Logger interface {
	Info(msg ...interface{}) LogItem
	Warn(msg ...interface{}) LogItem
	Error(msg ...interface{}) LogItem
	Debug(msg ...interface{}) LogItem
	WithAction(LogAction) Logger
	WithModule(string) Logger
}

type LogItem interface {
	Bytes() ([]byte, error)
	String() (string, error)
}
