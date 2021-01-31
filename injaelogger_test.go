package injaelogger

import (
	zerolog "github.com/rs/zerolog/log"
	"injaelogger/log"
	"strconv"
	"testing"
)

func TestLog(t *testing.T) {
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
