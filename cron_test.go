package golang_cron

import (
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	cron, err := New(time.Second)
	if err != nil {
		t.Fatal(err)
	}

	defer cron.Stop()

	err = cron.UnregisterJob("job")
	if err == nil {
		t.Fatal("Cron doesn't return error on unexpected job name deleting")
	}

	err = cron.DelSchedule(123)
	if err == nil {
		t.Fatal("Cron doesn't return error on unexpected job schedule deleting")
	}

	var counter int
	cron.RegisterJob("job", func() {
		counter++
	})

	scheduleId := cron.AddSchedule("job", func(tm time.Time) (ok bool) {
		return tm.Second()%2 == 0
	})

	time.Sleep(3 * time.Second)

	if counter != 1 {
		t.Fatal("Schedule works unexpectedly")
	}

	err = cron.DelSchedule(scheduleId)
	if err != nil {
		t.Fatal(err)
	}

	err = cron.UnregisterJob("job")
	if err != nil {
		t.Fatal(err)
	}
}
