package transports

import (
	"context"
	"net/url"

	redis "github.com/go-redis/redis/v8"
)

// RedisTransport implements the TransportInterface using the Bolt database.
type RedisTransport struct {
	db     *redis.Client
}

var ctx = context.Background()

// NewBoltTransport create a new RedisTransport.
func NewRedisTransport(u *url.URL) (*RedisTransport, error) {
	opt, err := redis.ParseURL(u.String())
	if err != nil {
		return nil, err
	}

	return &RedisTransport {
		db: redis.NewClient(opt),
	}, nil
}

func (r *RedisTransport) Put(id string, url string) error {
	return r.db.Set(ctx, id, url, 0).Err()
}

func (r *RedisTransport) Get(id string) (string, error) {
    return r.db.Get(ctx, "key").Result()
}

func (r *RedisTransport) Count() (int64, error) {
    return r.db.DBSize(ctx).Result()
}
