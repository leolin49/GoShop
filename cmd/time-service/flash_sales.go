package main

import (
	"goshop/pkg/util"
	"time"

	"github.com/golang/glog"
)

var (
	FlashKey             string
	FlashSalePeriod      = 1 * time.Hour
	FlashSalePrepare     = 10 * time.Minute
	FlashSaleDuration    = 5 * time.Minute
	FlashCacheWarmUpChan = make(chan time.Time)
	FlashBeginChan       = make(chan time.Time)
	FlashEndChan         = make(chan time.Time)
)

func startAllTicker() {
	go flashCacheWarmUpTicker()
	go flashSaleTicker()
	go flashSaleEndTicker()
}

func flashCacheWarmUpTicker() {
	if t := util.TimeToNextHour(); t > FlashSalePrepare {
		glog.Infoln(t, FlashSalePrepare, t-FlashSalePrepare)
		time.Sleep(t - FlashSalePrepare)
	}
	glog.Infoln("1111", time.Now())
	prepareTicker := time.NewTicker(FlashSalePeriod)
	glog.Infoln("2222", time.Now())
	defer prepareTicker.Stop()

	glog.Infoln("3333", time.Now())
	FlashCacheWarmUpChan <- time.Now()
	glog.Infoln("4444", time.Now())
	for {
		glog.Infoln("5555", time.Now())
		select {
		case <-prepareTicker.C:
			FlashCacheWarmUpChan <- time.Now()
		}
	}
}

func flashSaleTicker() {
	time.Sleep(util.TimeToNextHour())
	beginTicker := time.NewTicker(FlashSalePeriod)
	defer beginTicker.Stop()

	FlashBeginChan <- time.Now()
	for {
		select {
		case <-beginTicker.C:
			FlashBeginChan <- time.Now()
		}
	}
}

func flashSaleEndTicker() {
	time.Sleep(util.TimeToNextHour() + FlashSaleDuration)
	endTicker := time.NewTicker(FlashSalePeriod)
	defer endTicker.Stop()

	FlashEndChan <- time.Now()
	for {
		select {
		case <-endTicker.C:
			FlashEndChan <- time.Now()
		}
	}
}
