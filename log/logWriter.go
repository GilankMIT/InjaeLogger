package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func WriteLog(l *Logger) {
	log.WithLevel(zerolog.Level(l.logLevel)).Msg(l.message)
}

func logFileWriter(logText string) error {
	currentLogFileName := getLogFile()

	//check log file name, if different, re open log file
	if currentLogFileName != inUseLogFileName {
		inUseLogFileName = currentLogFileName
		f, err := os.OpenFile(LogFolder+"/"+currentLogFileName,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		currentLogFile = f
	}
	pushLog(logText)
	return nil
}

var logWriterPool = make(chan string, 50)

func pushLog(logText string) {
	logWriterPool <- logText

}

func logWriterPoolListen() {
	for {
		select {
		case logText := <-logWriterPool:
			if _, err := currentLogFile.WriteString(logText + "\n"); err != nil {
				log.Error().Msg(err.Error())
			}
		}
	}
}
