package ktimer

import (
	"github.com/robfig/cron"
	"math/rand"
	"sync"
	"time"
)

type KTimerBlock func(t *KTimer)

var timerDict map[int]*KTimer = map[int]*KTimer{}
var timerDictLock sync.Mutex = sync.Mutex{}

type KTimer struct {
	running   bool
	invalid   bool
	key       int
	fireTimes int
	block     func(t *KTimer)
	UserInfo  interface{}
	//ticker
	ticker      *time.Ticker
	stopChan    chan bool
	interval    time.Duration
	repeatTimes int
	//corn
	cron *cron.Cron
	spec string
}

func retainTimer(t *KTimer) {
	timerDictLock.Lock()
	if timerDict[t.key] == t {
		timerDictLock.Unlock()
		return
	}
	for {
		key := rand.Int()
		if timerDict[key] == nil {
			t.key = key
			timerDict[key] = t
			timerDictLock.Unlock()
			return
		}
	}
}

func releaseTimer(t *KTimer) {
	timerDictLock.Lock()
	timerDict[t.key] = nil
	timerDictLock.Unlock()
}

func (t *KTimer) run() {
	t.ticker = time.NewTicker(t.interval)
	for {
		select {
		case <-t.ticker.C:
			t.fireTimes++
			t.block(t)
			if t.repeatTimes != 0 && t.fireTimes >= t.repeatTimes {
				t.Stop(true)
			}
		case <-t.stopChan:
			t.ticker.Stop()
			return
		}
	}
}

func NewTimer(interval float64, repeatTimes int, block KTimerBlock) *KTimer {
	if block == nil {
		return nil
	}

	timer := &KTimer{interval: time.Duration(interval * float64(time.Second)), repeatTimes: repeatTimes, block: block, stopChan: make(chan bool, 1)}
	return timer
}

func NewCronTimer(spec string, block KTimerBlock) *KTimer {
	if block == nil {
		return nil
	}
	timer := &KTimer{spec: spec, block: block}
	timer.cron = cron.New()
	timer.cron.AddFunc(spec, func() {
		timer.fireTimes++
		block(timer)
	})
	return timer
}

func (t *KTimer) Start() {
	if !t.running && !t.invalid {
		retainTimer(t)
		t.running = true
		if t.cron != nil {
			t.cron.Start()
		} else {
			go t.run()
		}
	}
}

func (t *KTimer) Stop(invalidate bool) {
	if t.running {
		if t.cron != nil {
			t.cron.Stop()
		} else {
			t.stopChan <- true
		}
		t.running = false
	}
	if invalidate && !t.invalid {
		t.invalid = true
		if t.stopChan != nil {
			close(t.stopChan)
		}
		releaseTimer(t)
	}
}

func (t *KTimer) IsRunning() bool {
	return t.running
}

func (t *KTimer) IsValid() bool {
	return !t.invalid
}

func (t *KTimer) FireTimes() int {
	return t.fireTimes
}

func (t *KTimer) Fire() {
	t.block(t)
}
