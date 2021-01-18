# Golang jobs scheduler

## Create new cron

```go
import cron "github.com/leprosus/golang-cron"

c := cron.New(time.Minute)
```

## Add job and schedule

```go
c.RegisterJob("job", func () {
//Do something here
})

scheduleId := c.AddSchedule("job", func (tm time.Time) (ok bool) {
// It fires an execution of job "job" every even second

return tm.Second()%2 == 0
})
```

## Remove schedule and job

```go
c.DelSchedule("job", scheduleId)

c.UnregisterJob("job")
```

## List all methods

* New(ticker) - creates new cron
* RegisterJob(name, func) - registers the new job with own name and function to execute when a schedule time will come
* AddSchedule(name, handler) - adds schedule for a job with name and time handler and returns schedule ID
* DelSchedule(scheduleId) - delete schedule by ID
* UnregisterJob(name) - unregisters a job and all schedules that were connected with its