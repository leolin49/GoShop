package main

import (
	"context"
	"errors"
	"fmt"
	stockpb "goshop/api/protobuf/stock"
	"goshop/models"
	errorcode "goshop/pkg/error"
	"strconv"

	"github.com/golang/glog"
)

// TODO
func (s *StockRpcService) FlashStock(ctx context.Context, req *stockpb.ReqFlashStock) (*stockpb.RspFlashStock, error) {
	productKey := fmt.Sprintf("product_flash:%d", req.ProductId)
	if exist, err := rdb.Exist(productKey); err != nil {
		return nil, err
	} else if !exist {
		return nil, errors.New(
			fmt.Sprintf("[StockServer] flash stock failed, the stock [%s] not exist\n", productKey),
		)
	}
	res, err := rdb.RunScript(`
		local product_key = KEYS[1]
		local flash_count = ARGV[1]
		local stock = tonumber(redis.call('GET', product_key))
		if stock < flash_count then
			return 0
		end
		redis.call('DECR', flash_count)
		return 1
	`, []string{productKey}, []interface{}{req.SubCount})
	if err != nil {
		return nil, err
	}
	success, _ := res.(int)
	if success == 0 {
		return &stockpb.RspFlashStock{ErrorCode: errorcode.FlashNoStock}, nil
	}
	return &stockpb.RspFlashStock{ErrorCode: errorcode.Ok}, nil
}

func (s *StockRpcService) FlashCacheWarmUp(ctx context.Context, req *stockpb.ReqFlashCacheWarmUp) (*stockpb.RspFlashCacheWarmUp, error) {
	stocks, err := models.NewStockQuery(db).GetAllStock()
	if err != nil {
		glog.Errorln("[StockServer] get all stock error: ", err)
		return nil, err
	}
	for _, st := range stocks {
		err = rdb.Set(fmt.Sprintf("product_flash:%d", st.ProductId), strconv.Itoa(int(st.Count)))
		if err != nil {
			glog.Errorf("[StockServer] Flash cache warm up failed [%v] on product [%v]\n", err, st)
			return nil, err
		}
	}
	return &stockpb.RspFlashCacheWarmUp{ErrorCode: errorcode.Ok}, nil
}

func (s *StockRpcService) FlashCacheClean(ctx context.Context, req *stockpb.ReqFlashCacheClear) (*stockpb.RspFlashCacheClear, error) {
	stocks, err := models.NewStockQuery(db).GetAllStock()
	if err != nil {
		glog.Errorln("[StockServer] get all stock error: ", err)
		return nil, err
	}
	for _, st := range stocks {
		key := fmt.Sprintf("product_flash:%d", st.ProductId)
		res, err := rdb.RunScript(`
			local cnt = tonumber(redis.call('get', KEYS[1]))
			redis.call('del', KEYS[1])
			return cnt
		`, []string{key}, []interface{}{})
		if err != nil {
			return nil, err
		}

		cnt, _ := res.(uint64)
		// write to mysql
		go func(productId uint32, count uint64) {
			if err := models.NewStockQuery(db).SetStock(st.ProductId, count); err != nil {
				glog.Errorf("[StockServer] write the cache [%d - %d] back to mysql failed: %v", productId, count, err)
				return
			}
		}(st.ProductId, cnt)
	}
	return &stockpb.RspFlashCacheClear{ErrorCode: errorcode.Ok}, nil
}
