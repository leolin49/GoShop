package main

import (
	"goshop/pkg/util"
	"time"
)

var (
	FlashSalePeriod      = 1 * time.Hour
	FlashSaleDuration    = 5 * time.Minute
	FlashCacheWarmUpChan = make(chan time.Time)
	FlashBeginChan       = make(chan time.Time)
	FlashEndChan         = make(chan time.Time)
)

func startAllTicker() {
	go flashCacheWarmUpTicker()
	go flashSaleTicker()
}

func flashCacheWarmUpTicker() {
	if util.TimeToNextHour() > 5*time.Minute {
		time.Sleep(util.TimeToNextHour() - 5*time.Minute)
	}
	prepareTicker := time.NewTicker(FlashSalePeriod)
	FlashCacheWarmUpChan <- time.Now()
	defer prepareTicker.Stop()

	for {
		select {
		case <-prepareTicker.C:
			FlashCacheWarmUpChan <- time.Now()
		}
	}
}

func flashSaleTicker() {
	time.Sleep(util.TimeToNextHour())
	beginTicker := time.NewTicker(FlashSalePeriod)
	FlashBeginChan <- time.Now()
	defer beginTicker.Stop()

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
	FlashEndChan <- time.Now()
	defer endTicker.Stop()

	for {
		select {
		case <-endTicker.C:
			FlashEndChan <- time.Now()
		}
	}
}
