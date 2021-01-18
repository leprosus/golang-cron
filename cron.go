package golang_cron

import (
	"fmt"
	"sync"
	"time"
)

type Cron struct {
	mx *sync.RWMutex

	ticker *time.Ticker

	ix int

	jobs      map[name]job
	schedules map[int]schedule

	ids   map[name][]int
	names map[int]name
}

type job func()

type name interface{}

type schedule func(tm time.Time) (ok bool)

func New(ticker time.Duration) (c *Cron, err error) {
	if ticker < time.Second {
		err = fmt.Errorf("the ticker duration is set too low: %v", ticker)

		return
	}

	c = &Cron{
		mx: &sync.RWMutex{},

		ticker: time.NewTicker(ticker),

		jobs:      map[name]job{},
		schedules: map[int]schedule{},

		ids:   map[name][]int{},
		names: map[int]name{},
	}

	go func() {
		for range c.ticker.C {
			c.tock()
		}
	}()

	return
}

func (c *Cron) tock() {
	c.mx.RLock()
	defer c.mx.RUnlock()

	var n name
	tm := time.Now()
	for id, s := range c.schedules {
		if s(tm) {
			n = c.names[id]
			c.jobs[n]()
		}
	}
}

func (c *Cron) RegisterJob(name interface{}, job func()) {
	c.mx.Lock()
	c.jobs[name] = job
	c.ids[name] = []int{}
	c.mx.Unlock()

	return
}

func (c *Cron) AddSchedule(name interface{}, handler func(tm time.Time) (ok bool)) (scheduleId int) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.ix++
	scheduleId = c.ix

	c.schedules[scheduleId] = handler

	ids := c.ids[name]
	ids = append(ids, scheduleId)
	c.ids[name] = ids

	c.names[scheduleId] = name

	return
}

func (c *Cron) DelSchedule(scheduleId int) (err error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	curName := c.names[scheduleId]

	_, ok := c.schedules[scheduleId]
	if !ok {
		err = fmt.Errorf("there isn't schedule with ID %d", scheduleId)

		return
	}

	delete(c.schedules, scheduleId)

	ids := c.ids[curName]
	for i, id := range ids {
		if id == scheduleId {
			c.ids[curName] = append(ids[:i], ids[i+1:]...)

			ok = true

			break
		}
	}

	if !ok {
		err = fmt.Errorf("there isn't schedule with ID %d into namespace %s", scheduleId, curName)

		return
	}

	delete(c.names, scheduleId)

	return
}

func (c *Cron) UnregisterJob(name interface{}) (err error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	_, ok := c.jobs[name]
	if !ok {
		err = fmt.Errorf("there isn't job with name %s", name)

		return
	}

	for _, id := range c.ids[name] {
		delete(c.schedules, id)
		delete(c.names, id)
	}

	delete(c.ids, name)

	return
}

func (c *Cron) Stop() {
	c.ticker.Stop()
}
