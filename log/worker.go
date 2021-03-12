package log

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	workerCheckerDuration = time.Second * 5
	logRotationDuration   = time.Hour * 24 * 30 //set log rotation to 30days default
)

func SetLogRotationDuration(duration time.Duration) error {
	if duration.Minutes() < 1 { //forbid if duration is less than one minute for performance purpose
		return errors.New("cannot set log time rotation to less than one minute")
	}

	logRotationDuration = duration

	refreshWorker()
	return nil
}

func refreshWorker() {

}

type ScheduledWorker struct {
	logRotationTicker *time.Ticker
	workerTicker      *time.Ticker
}

func NewScheduledWorker() ScheduledWorker {
	return ScheduledWorker{
		logRotationTicker: time.NewTicker(logRotationDuration),
		workerTicker:      time.NewTicker(workerCheckerDuration),
	}
}

func WithLogRotationWorker() {
	worker := NewScheduledWorker()
	worker.logRotationWorker()
}

func (worker *ScheduledWorker) logRotationWorker() {
	fmt.Println("Log rotation worker running")
	for {
		select {
		case <-worker.workerTicker.C:
			fmt.Println("Log rotation worker running on : " + time.Now().String())
			worker.doLogRotation()
		}
	}

}

func (worker *ScheduledWorker) doLogRotation() {
	//get the date of the log rotation and check if the file is exist
	//logRotationDate := time.Now().Add(logRotationDuration)

	//check log folder if there is any file exist
	//get rotated log deleted time
	rotatedLogDeleteTime := time.Now().Add(-logRotationDuration)

	root := LogFolder
	deletedLogFiles := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		//get log date time from filename
		fileNameOnlyLog := ""
		pathSeparated := strings.Split(path, "/")
		fileNameOnlyLog = pathSeparated[len(pathSeparated)-1]

		logTime, err := worker.getLogFileDateTime(fileNameOnlyLog)
		if err != nil {
			return nil
		}

		if rotatedLogDeleteTime.Before(logTime) {
			deletedLogFiles = append(deletedLogFiles, path)
		}
		return nil
	})
	if err != nil {
		log.Error().Msg("doLogRotaion err " + err.Error())
		return
	}

	err = worker.batchDeleteLogFile(deletedLogFiles)
	if err != nil {
		log.Error().Msg("doLogRotation err " + err.Error())
	}

	return
}

func (worker *ScheduledWorker) getLogFileDateTime(fileName string) (time.Time, error) {
	fileNameSeparated := strings.Split(fileName, "_")
	if len(fileNameSeparated) < 2 {
		return time.Time{}, nil
	}

	//check extension
	fileNameDate := ""
	fileNameDateSeparated := strings.Split(fileNameSeparated[1], ".")
	if len(fileNameDateSeparated) < 2 {
		fileNameDate = fileNameSeparated[1]
	} else {
		fileNameDate = fileNameSeparated[1][:strings.IndexByte(fileNameSeparated[1], '.')]
	}

	logTime, err := time.Parse("2006-01-02", fileNameDate)
	if err != nil {
		return time.Time{}, err
	}

	return logTime, nil
}

func (worker *ScheduledWorker) batchDeleteLogFile(logFiles []string) error {
	for _, file := range logFiles {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}
