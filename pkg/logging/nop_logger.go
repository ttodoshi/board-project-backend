package logging

var noOperationLogger *NoOperationLogger

type NoOperationLogger struct {
}

func GetNoOperationLogger() Logger {
	return noOperationLogger
}

func init() {
	noOperationLogger = &NoOperationLogger{}
}

func (l *NoOperationLogger) Print(_ ...interface{}) {
}

func (l *NoOperationLogger) Printf(_ string, _ ...interface{}) {
}

func (l *NoOperationLogger) Trace(_ ...interface{}) {
}

func (l *NoOperationLogger) Tracef(_ string, _ ...interface{}) {
}

func (l *NoOperationLogger) Debug(_ ...interface{}) {
}

func (l *NoOperationLogger) Debugf(_ string, _ ...interface{}) {
}

func (l *NoOperationLogger) Info(_ ...interface{}) {
}

func (l *NoOperationLogger) Infof(_ string, _ ...interface{}) {
}

func (l *NoOperationLogger) Warn(_ ...interface{}) {
}

func (l *NoOperationLogger) Warnf(_ string, _ ...interface{}) {
}

func (l *NoOperationLogger) Error(_ ...interface{}) {
}

func (l *NoOperationLogger) Errorf(_ string, _ ...interface{}) {
}

func (l *NoOperationLogger) Fatal(_ ...interface{}) {
}

func (l *NoOperationLogger) Fatalf(_ string, _ ...interface{}) {
}

func (l *NoOperationLogger) Panic(_ ...interface{}) {
}

func (l *NoOperationLogger) Panicf(_ string, _ ...interface{}) {
}
