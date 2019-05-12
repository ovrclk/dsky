package dsky

type shellLogger struct {
}

func (j *shellLogger) Info(msg ...interface{}) LogItem {
	return nil
}

func (j *shellLogger) Warn(msg ...interface{}) LogItem {
	return nil
}

func (j *shellLogger) Error(msg ...interface{}) LogItem {
	return nil
}

func (j *shellLogger) Debug(msg ...interface{}) LogItem {
	return nil
}

func (j *shellLogger) WithModule(string) Logger {
	return j
}

type shellLogItem struct{}

func (j *shellLogItem) Bytes() ([]byte, error) {
	return nil, nil
}

func (j *shellLogItem) String() (string, error) {
	return "", nil
}
