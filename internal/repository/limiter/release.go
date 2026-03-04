package limiter

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var releaseScript = redis.NewScript(`
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

func (r *Repository) Release(ctx context.Context, key string) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("redis client is nil")
	}

	if key == "" {
		return fmt.Errorf("empty redis key")
	}

	_, err := releaseScript.Run(ctx, r.db, []string{key}).Int64()
	return err
}
