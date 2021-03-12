package injaelogger

import (
	"fmt"
	"github.com/GilankMIT/InjaeLogger/log"
	zerolog "github.com/rs/zerolog/log"
	"strconv"
	"testing"
	"time"
)

func TestLogDebug(t *testing.T) {
	log.AppName = "injaelogger"
	log.LogFolder = "./example_log"
	log.WriteLogToFile(true)

	waitChan := make(chan int8)
	go func() {
		for i := 0; i < 1000; i++ {
			log.Debug().Msg("This is a log number " + strconv.Itoa(i))
		}
		waitChan <- 1
	}()

	<-waitChan
}

func TestLogInfo(t *testing.T) {
	log.AppName = "injaelogger"
	log.LogFolder = "./example_log"
	log.WriteLogToFile(true)

	waitChan := make(chan int8)
	go func() {
		for i := 0; i < 1000; i++ {
			log.Info().Msg("This is a log number " + strconv.Itoa(i))
		}
		waitChan <- 1
	}()

	<-waitChan
}

func TestLogWarn(t *testing.T) {
	log.AppName = "injaelogger"
	log.LogFolder = "./example_log"
	log.WriteLogToFile(true)

	waitChan := make(chan int8)
	go func() {
		for i := 0; i < 1000; i++ {
			log.Warn().Msg("This is a log number " + strconv.Itoa(i))
		}
		waitChan <- 1
	}()

	<-waitChan
}

func TestLogError(t *testing.T) {
	log.AppName = "injaelogger"
	log.LogFolder = "./example_log"
	log.WriteLogToFile(true)

	waitChan := make(chan int8)
	go func() {
		for i := 0; i < 1000; i++ {
			log.Error().Msg("This is a log number " + strconv.Itoa(i))
		}
		waitChan <- 1
	}()

	<-waitChan
}

func BenchmarkLog(b *testing.B) {
	log.AppName = "injaelogger"
	log.LogFolder = "./example_log"
	log.WriteLogToFile(true)

	for i := 0; i < b.N; i++ {
		log.Info().Msg("This is a log")
	}
}

func BenchmarkZerolog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		zerolog.Info().Msg("")
	}
}

func TestLogRotation(t *testing.T) {
	log.AppName = "InjaeLoggerTest"
	log.LogFolder = "./example_log"
	err := log.SetLogRotationDuration(time.Hour * 24)
	if err != nil {
		t.Error(err)
		return
	}

	log.WriteLogToFile(true)
	for i := 0; i < 100; i++ {
		log.Info().Msg("This is a log")
	}
	log.WithLogRotationWorker()
}

func TestLogWriteLevel(t *testing.T) {
	log.AppName = "InjaeLoggerWriteFile"
	log.LogFolder = "./log_folder"
	log.WriteLogToFile(true)
	log.SetLogRotationDuration(time.Hour * 24 * 90)
	for i := 0; i < 100; i++ {
		log.Write("This is a log to file, no std out " + strconv.Itoa(i))
	}

	//do sleep to ensure log writer finish the job
	sleeper := time.NewTimer(time.Second * 10)
	<-sleeper.C

	fmt.Println("Logger finished working")
}
