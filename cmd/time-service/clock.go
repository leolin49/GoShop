package main

import (
	"context"
	"fmt"
	stockpb "goshop/api/protobuf/stock"

	"github.com/golang/glog"
)

func clockTrigger() {
	for {
		select {
		// Flash sales prepare. 23:50:00
		// 1. warm up the cache. (tran the data to redis).
		case t := <-FlashCacheWarmUpChan:
			glog.Infof("[TimeServer] flash cache warm up start in %v\n", t)
			if _, err := StockClient().FlashCacheWarmUp(context.Background(), &stockpb.ReqFlashCacheWarmUp{}); err != nil {
				glog.Errorf("[TimeServer] flash cache warm up failed: %v\n", err)
			}
		// Flash sales start. 00:00:00
		// 1. set the flash:Hour key in redis as the flag.
		case t := <-FlashBeginChan:
			glog.Infof("[TimeServer] flash sales start in %v", t)
			FlashKey = fmt.Sprintf("flash_sales:%d", t.Unix())
			if err := rdb.SetInt(FlashKey, 1); err != nil {
				glog.Errorf("[TimeServer] flash sales start failed: %v\n", err)
			}
		// Flash sales end. 00:05:00
		// 1. clean the cache.
		// 2. tran the cache data to mysql.
		case t := <-FlashEndChan:
			glog.Infof("[TimeServer] flash sales end in %v\n", t)
			if err := rdb.Del(FlashKey); err != nil {
				glog.Errorf("[TimeServer] flash sales end failed: %v\n", err)
			}
			if _, err := StockClient().FlashCacheClear(context.Background(), &stockpb.ReqFlashCacheClear{}); err != nil {
				glog.Errorf("[TimeServer] flash cache clean failed: %v\n", err)
			}
		}
	}
}
