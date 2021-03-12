package log

import "testing"

func TestGetLogFileDateTime(t *testing.T) {
	worker := ScheduledWorker{}
	logTime, err := worker.getLogFileDateTime("AppName_2021-02-23")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(logTime)
}
