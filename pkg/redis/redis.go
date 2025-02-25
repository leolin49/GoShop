package redis

import (
	"context"
	"errors"
	"fmt"
	"goshop/configs"
	"goshop/pkg/util"
	"strconv"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)
type Rdb struct {
	ctx context.Context
	db *redis.Client
}

func NewRedisClient(cfg *configs.RedisConfig) (*Rdb, error) {
	rdb := &Rdb{
		db:	redis.NewClient(&redis.Options{
			Addr:	fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password: cfg.Password,	// no password set
			DB: cfg.Db,	// use default DB
			Protocol: cfg.Protocol,	// specify 2 for RESP 2, or 3 for RESP 3.
		}),
	}
	if rdb.db == nil {
		return nil, errors.New("redis new client failed")
	}
	if !rdb.Ping() {
		return nil, errors.New("redis ping failed")
	}
	return rdb, nil
}

// Ping
func (r *Rdb) Ping() bool {
	_, err := r.db.Ping(r.ctx).Result()
	return err == nil
}

// Exist
func (r *Rdb) Exist(k string) (bool, error) {
	return r.Exist(k)
}

func (r *Rdb) Del(k string) error {
	return r.Del(k)
}

// String
func (r *Rdb) Set(k, v string) error {
	return r.db.Set(r.ctx, k, v, 0).Err()
}

func (r *Rdb) Get(k string) (string, error) {
	return r.db.Get(r.ctx, k).Result()
}

func (r *Rdb) SetInt(k string, v int) error {
	val := strconv.Itoa(v)
	return r.db.Set(r.ctx, k, val, 0).Err()
}

func (r *Rdb) GetInt(k string) (int, error) {
	return r.db.Get(r.ctx, k).Int()
}

func (r *Rdb) SetProto(k string, v proto.Message) error {
	d, err := util.Serialize(v)
	if err != nil {
		return err
	}
	return r.Set(k, string(d))
}

func (r *Rdb) GetProto(k string, v proto.Message) (bool, error) {
	d, err := r.Get(k)
	if err != nil {
		if err = util.Deserialize([]byte(d), v); err == nil {
			return true, nil
		}
	} else if err != redis.Nil {
		// redis error
		return false, err	
	}
	// cache not hit
	return false, nil
}


