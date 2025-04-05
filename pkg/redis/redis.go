package redis

import (
	"context"
	"errors"
	"fmt"
	"goshop/configs"
	"goshop/pkg/util"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

type IRdb interface {
	Ping() bool
	Exist(k string) (bool, error)
	Del(k string) error
	Set(k, v string) error
	Get(k string) (string, error)
	SetInt(k string, v int) error
	GetInt(k string) (int, error)
	SetProto(k string, v proto.Message) error
	GetProto(k string, v proto.Message) (bool, error)
	RunScript(src string, keys []string, args []interface{}) (int, error)
	Exec(args ...interface{}) *redis.Cmd
	Lock(k string) (bool, error)
	TryLock(k string) (bool, error)
	LockFunc(k string, f func() (interface{}, error)) (interface{}, error)
	UnLock(k string) bool
}

type Rdb struct {
	ctx context.Context
	db  *redis.Client
}

func NewRedisClient(cfg *configs.RedisConfig) (IRdb, error) {
	rdb := &Rdb{
		ctx: context.Background(),
		db: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password: cfg.Password, // no password set
			DB:       cfg.Db,       // use default DB
			Protocol: cfg.Protocol, // specify 2 for RESP 2, or 3 for RESP 3.
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
	v, err := r.db.Exists(r.ctx, k).Result()
	if err != nil {
		return false, err
	}
	return v > 0, nil
}

func (r *Rdb) Del(k string) error {
	return r.db.Del(r.ctx, k).Err()
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

	switch {
	case err == redis.Nil:
		return false, nil
	case err != nil:
		return false, err
	}
	if err = util.Deserialize([]byte(d), v); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Rdb) RunScript(src string, keys []string, args []interface{}) (int, error) {
	srcript := redis.NewScript(src)
	return srcript.Run(r.ctx, r.db, keys, args...).Int()
}

func (r *Rdb) Exec(args ...interface{}) *redis.Cmd {
	return r.db.Do(r.ctx, args)
}

func (r *Rdb) Lock(k string) (bool, error) {
	res, err := r.Exec("SET", k, 1, "NX", "PX", 1000).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func (r *Rdb) TryLock(k string) (bool, error) {
	maxRetry := 5
	interval := time.Second
	for range maxRetry {
		if ok, err := r.Lock(k); ok && err == nil {
			return true, nil
		}
		time.Sleep(interval)
		interval = time.Duration(float64(interval) * 1.5)
	}
	return false, errors.New("redis try lock failed")
}

func (r *Rdb) LockFunc(k string, f func() (interface{}, error)) (interface{}, error) {
	locked, err := r.TryLock(k)
	defer r.UnLock(k)
	if err != nil || !locked {
		return nil, err
	}
	// if !locked {
	// 	return nil, errors.New(fmt.Sprintf("lock [%s] is locked by others.", k))
	// }
	return f()
}

func (r *Rdb) UnLock(k string) bool {
	err := r.Del(k)
	return err == nil
}
