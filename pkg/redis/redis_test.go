package redis

import (
	cartpb "goshop/api/protobuf/cart"
	"goshop/pkg/util"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockIRdb(ctrl)

	mockRedis.EXPECT().
		Ping().
		Return(true)

	assert.True(t, mockRedis.Ping())
}

func TestSetAndGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockIRdb(ctrl)

	mockRedis.EXPECT().
		Set("key", "value").
		Return(nil)
	err := mockRedis.Set("key", "value")
	assert.NoError(t, err)

	mockRedis.EXPECT().
		Get("key").
		Return("value", nil)
	val, err := mockRedis.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)

	// Get Failed
	mockRedis.EXPECT().
		Get("nonexistent").
		Return("", redis.Nil)
	_, err = mockRedis.Get("nonexistent")
	assert.Equal(t, redis.Nil, err)
}

func TestSetIntAndGetInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRdb := NewMockIRdb(ctrl)

	mockRdb.EXPECT().
		SetInt("key", 123).
		Return(nil)

	err := mockRdb.SetInt("key", 123)
	assert.NoError(t, err)

	mockRdb.EXPECT().
		GetInt("key").
		Return(123, nil)

	val, err := mockRdb.GetInt("key")
	assert.NoError(t, err)
	assert.Equal(t, 123, val)

	mockRdb.EXPECT().
		GetInt("nonexistent").
		Return(0, redis.Nil)

	_, err = mockRdb.GetInt("nonexistent")
	assert.Equal(t, redis.Nil, err)
}

func TestSetProtoAndGetProto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRdb := NewMockIRdb(ctrl)

	var items []*cartpb.CartItem
	items = append(items, &cartpb.CartItem{ProductId: 1, Quantity: 1})
	items = append(items, &cartpb.CartItem{ProductId: 12, Quantity: 10})
	items = append(items, &cartpb.CartItem{ProductId: 123, Quantity: 100})
	msg := &cartpb.RspGetCart{
		ErrorCode: 0,
		Cart: &cartpb.Cart{
			UserId: 1,
			Items:  items,
		},
	}
	serialized, _ := util.Serialize(msg)

	mockRdb.EXPECT().
		SetProto("key", msg).
		Return(nil)
	err := mockRdb.SetProto("key", msg)
	assert.NoError(t, err)

	mockRdb.EXPECT().
		GetProto("key", gomock.Any()).
		DoAndReturn(func(k string, v proto.Message) (bool, error) {
			if err := util.Deserialize([]byte(serialized), v); err != nil {
				return false, err
			}
			return true, nil
		})

	var result cartpb.RspGetCart
	exists, err := mockRdb.GetProto("key", &result)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, int32(0), result.ErrorCode)
	assert.Equal(t, uint32(1), result.Cart.UserId)
	assert.Equal(t, 3, len(result.Cart.Items))
	for i, item := range result.Cart.Items {
		if i == 0 {
			assert.Equal(t, uint32(1), item.ProductId)
			assert.Equal(t, int32(1), item.Quantity)
		} else if i == 1 {
			assert.Equal(t, uint32(12), item.ProductId)
			assert.Equal(t, int32(10), item.Quantity)
		} else if i == 2 {
			assert.Equal(t, uint32(123), item.ProductId)
			assert.Equal(t, int32(100), item.Quantity)
		}
	}

	mockRdb.EXPECT().
		GetProto("nonexistent", gomock.Any()).
		Return(false, nil)

	exists, err = mockRdb.GetProto("nonexistent", &result)
	assert.NoError(t, err)
	assert.False(t, exists)
}
