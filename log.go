package dsky

type logModeType string

const (
	logModeTypeInfo  logModeType = "info"
	logModeTypeDebug             = "debug"
	logModeTypeWarn              = "warn"
	logModeTypeError             = "error"
)

type Logger interface {
	Info(msg ...interface{}) LogItem
	Warn(msg ...interface{}) LogItem
	Error(msg ...interface{}) LogItem
	Debug(msg ...interface{}) LogItem
	WithModule(string) Logger
}

type LogItem interface {
	Bytes() ([]byte, error)
	String() (string, error)
}
