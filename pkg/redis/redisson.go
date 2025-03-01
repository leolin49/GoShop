package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const RedissonUnlockMessage = 0

func TryLock(ctx context.Context, rdb *redis.Client, lockKey string, lockId string) (bool, error) {
	var tryLockScript = redis.NewScript(`
		local lock_key = KEYS[1]
		local expire_time_ms = ARGV[1]	
		local lock_id = ARGV[2]

		-- lock not exist, add lock and set lockcnt to 1.
		if (redis.call('exists', lock_key) == 0) then
			redis.call('hset', lock_key, lock_id, 1)
			redis.call('pexpire', lock_key, expire_time_ms)
			return nil
		end

		-- lock exist and lock_id match, increase the lockcnt.
		if (redis.call('hexists', lock_key, lock_id) == 1) then
			redis.call('hincrby', lock_key, lock_id, 1)
			redis.call('pexpire', lock_key, expire_time_ms)
			return nil
		end

		-- locked by other thread, return the lock expire time.
		return redis.call('pttl', lock_key)
		`)
	keys := []string{lockKey}
	args := []interface{}{10 * 1000, lockId}
	res, err := tryLockScript.Run(ctx, rdb, keys, args...).Result()
	if err != nil {
		return false, err
	}
	return res == nil, nil
}

func UnLock(ctx context.Context, rdb *redis.Client, lockKey string, lockId string) (bool, error) {
	var unLockScript = redis.NewScript(`
		local lock_key = KEYS[1]
		local pub_key = KEYS[2]
		local unlock_msg = ARGV[1]
		local expire_time_ms = ARGV[2]	
		local lock_id = ARGV[3]

		-- if lock not exist, publish the unlock message to the 'pub_key' channel.
		if (redis.call('exists', lock_key) == 0) then	
			redis.call('publish', pub_key, unlock_msg)
			return 1
		end

		-- lock exist and lock_id not match, locked by other thread.
		if (redis.call('hexists', lock_key, lock_id) == 0) then
			return 2
		end

		-- lock exist and lock_id match
		local counter = redis.call('hincrby', lock_key, lock_id, -1)
		if (counter > 0) then
			-- other same id thread locked, increase some lock time for it.
			redis.call('pexpire', lock_key, expire_time_ms)
			return 0
		else
			-- counter is sub to zero, no any theard keep the lock,
			-- 1. delete the lock_key
			-- 2. publish the unlock message
			redis.call('del', lock_key)
			redis.call('publish', pub_key, unlock_msg)
			return 1
		end

		return nil
		`)
	keys := []string{lockKey, lockKey}
	args := []interface{}{RedissonUnlockMessage, 10 * 1000, lockId}
	res, err := unLockScript.Run(ctx, rdb, keys, args...).Result()
	if err != nil {
		return false, err
	}
	status, _ := res.(int)
	switch status {
	case 0:
		return true, nil
	case 1:
		return true, nil
	case 2:
		return false, nil
	}
	return false, nil
}
