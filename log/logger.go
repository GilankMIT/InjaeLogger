package log

/*
	InjaeLogger 2020
	Created by Gilang Prambudi (itgilangprambudi@gmail.com)

	A custom log wrapper for zerolog.
	Support:
	1. Standardize log file write
	2. Limiting log severity level
	3. Log rotation handler (WIP)
	4. Periodic push (batch upload) for logger (WIP)
	5. Instant log push (WIP)
*/

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	LogFolder         = "./log" //default log folder
	currentLogFile    *os.File
	inUseLogFileName  string
	RotationDuration  = time.Hour * 24 * 30 //default rotation day 30 day
	AppName           = ""
	writeLog          = false
	fileRuntimeCaller = 3 //default is 3
)

type Level uint8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// WriteLevel define write level log (no output in std.out, only write to file)
	WriteLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	case WriteLevel:
		return "write"
	}
	return ""
}

func WriteLogToFile(enable bool) {
	//ensure log folder existence (if log is defined)
	if enable {
		writeLog = enable
		err := createLogFolderIfNotExist()
		if err != nil {
			panic(err.Error())
		}

		go logWriterPoolListen()
	}
}

func createLogFolderIfNotExist() error {
	folderExist, err := checkDirExistence()
	if err != nil {
		return err
	}
	if !folderExist {
		err = os.MkdirAll(LogFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return err
}

func checkDirExistence() (bool, error) {
	_, err := os.Stat(LogFolder)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type Logger struct {
	message         string
	completeMessage map[string]string
	logLevel        Level
	fileName        string
	funcName        string
}

func (l *Logger) Str(fieldName, value string) *Logger {
	l.completeMessage[fieldName] = value
	return l
}

func (l *Logger) Msg(message string) {
	l.message = message

	//print log by building log message
	l.completeMessage["message"] = message
	jsonMessage := l.buildMessage()
	l.sendLog(jsonMessage)
}

func (l *Logger) buildMessage() string {
	jsonMessage, _ := json.Marshal(&l.completeMessage)
	return string(jsonMessage)
}

type logTextParam struct {
	fileName     string
	functionName string
	logLevel     string
	message      string
}

func (l *Logger) sendLog(message string) {
	if writeLog {
		err := logFileWriter(message)
		if err != nil {
			log.Println(err)
		}
	}

	if l.logLevel == WriteLevel {
		return
	}
	OutputLog(l)
}

func getLogTrace() (fileName string, funcName string) {

	//runtime.Caller skip is 3 because it is called from logging function -> newLogger() -> getLogTrace()
	pc, file, line, ok := runtime.Caller(fileRuntimeCaller)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]

	return filename, fn
}

func getLogFile() string {
	if AppName != "" {
		return AppName + "_" + time.Now().Format("2006-01-02") + ".log"
	}
	return "log_" + time.Now().Format("2006-01-02") + ".log"
}

//SetCaller is used in case there is other dependency that want to integrate InjaeLogger
func SetCaller(caller int) {
	fileRuntimeCaller = caller
}
