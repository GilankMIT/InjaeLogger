package log

func newLogger(logLevel Level) *Logger {
	fileName, funcName := getLogTrace()
	l := Logger{}
	l.logLevel = logLevel
	l.fileName = fileName
	l.funcName = funcName
	l.completeMessage = make(map[string]string)
	l.completeMessage["file"] = l.fileName
	l.completeMessage["app"] = AppName
	l.completeMessage["function"] = l.funcName
	l.completeMessage["level"] = logLevel.String()
	return &l
}

func Debug() *Logger {
	l := newLogger(DebugLevel)
	return l
}

func Info() *Logger {
	l := newLogger(InfoLevel)
	return l
}

func Warn() *Logger {
	l := newLogger(WarnLevel)
	return l
}

func Error() *Logger {
	l := newLogger(ErrorLevel)
	return l
}

func Write(message string) *Logger {
	l := newLogger(WriteLevel)
	l.completeMessage["message"] = message

	jsonMessage := l.buildMessage()
	l.sendLog(jsonMessage)
	return l
}
