package redis

import (
	"context"
	"fmt"
	"time"

	r "github.com/redis/go-redis/v9"
)

var reserveScript = r.NewScript(`
local current = tonumber(redis.call("GET", KEYS[1]) or "0")
local limit = tonumber(ARGV[1])

if current >= limit then
	return 0
end

redis.call("INCR", KEYS[1])

local ttl = tonumber(ARGV[2])
if ttl ~= nil and ttl > 0 then
	redis.call("EXPIRE", KEYS[1], ttl)
end

return 1
`)

var releaseScript = r.NewScript(`
if redis.call("EXISTS", KEYS[1]) == 0 then
	return 0
end

local nextValue = tonumber(redis.call("DECR", KEYS[1]))
if nextValue <= 0 then
	redis.call("DEL", KEYS[1])
	return 0
end

return nextValue
`)

func (c *Client) Reserve(
	ctx context.Context,
	key string,
	limit int64,
	ttl time.Duration,
) (bool, error) {
	if c == nil || c.client == nil {
		return false, fmt.Errorf("redis client is nil")
	}

	if key == "" {
		return false, fmt.Errorf("empty redis key")
	}

	if limit <= 0 {
		return false, fmt.Errorf("limit must be positive")
	}

	ttlSeconds := int64(0)
	if ttl > 0 {
		ttlSeconds = int64(ttl.Seconds())
		if ttlSeconds == 0 {
			ttlSeconds = 1
		}
	}

	result, err := reserveScript.Run(ctx, c.client, []string{key}, limit, ttlSeconds).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (c *Client) Release(ctx context.Context, key string) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("redis client is nil")
	}

	if key == "" {
		return fmt.Errorf("empty redis key")
	}

	_, err := releaseScript.Run(ctx, c.client, []string{key}).Int64()
	return err
}
