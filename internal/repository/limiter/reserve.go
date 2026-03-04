package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var reserveScript = redis.NewScript(`
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

func (r *Repository) Reserve(
	ctx context.Context,
	key string,
	limit int64,
	ttl time.Duration,
) (bool, error) {
	if r == nil || r.db == nil {
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

	result, err := reserveScript.Run(ctx, r.db, []string{key}, limit, ttlSeconds).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}
