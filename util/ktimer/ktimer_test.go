package ktMicro

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	wg := sync.WaitGroup{}
	repeatTimes := 5
	wg.Add(repeatTimes)
	NewTimer(0.333, repeatTimes, func(t *KTimer) {
		fmt.Printf("fire:%d %v userInfo:%v\n", t.FireTimes(), time.Now(), t.UserInfo)
		t.UserInfo = "data" + strconv.Itoa(t.FireTimes())
		wg.Done()
	}).Start()
	wg.Wait()
}

func TestTickerStop(t *testing.T) {
	wg := sync.WaitGroup{}
	repeatTimes := 5
	interval := 0.333
	wg.Add(repeatTimes)
	timer := NewTimer(interval, repeatTimes, func(t *KTimer) {
		fmt.Printf("fire:%d %v userInfo:%v\n", t.FireTimes(), time.Now(), t.UserInfo)
		t.UserInfo = "data" + strconv.Itoa(t.FireTimes())
		wg.Done()
	})
	timer.UserInfo = "SomeData"
	timer.Start()
	time.Sleep(time.Duration(float64(repeatTimes)*interval*float64(time.Second) - interval*0.5))
	timer.Stop(true)
	wg.Done()
	wg.Wait()
}

//func TestCron(t *testing.T){
//
//}
