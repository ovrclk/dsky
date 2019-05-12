package dsky

type jsonLogger struct {
}

func (j *jsonLogger) Info(msg ...interface{}) LogItem {
	return nil
}

func (j *jsonLogger) Warn(msg ...interface{}) LogItem {
	return nil
}

func (j *jsonLogger) Error(msg ...interface{}) LogItem {
	return nil
}

func (j *jsonLogger) Debug(msg ...interface{}) LogItem {
	return nil
}

func (j *jsonLogger) WithModule(string) Logger {
	return j
}

type jsonLogItem struct{}

func (j *jsonLogItem) Bytes() ([]byte, error) {
	return nil, nil
}

func (j *jsonLogItem) String() (string, error) {
	return "", nil
}
