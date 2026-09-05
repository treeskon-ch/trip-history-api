package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type locationCache struct {
	client *redis.Client
}

func NewLocationCache(client *redis.Client) ports.LocationCache {
	return &locationCache{
		client: client,
	}
}

func (c *locationCache) UpdateLatestLocation(ctx context.Context, tripID string, loc domain.LocationUpdate) error {
	redisKey := fmt.Sprintf("tracking:live:%s", tripID)
	redisData := fmt.Sprintf(`{"lat": %f, "lng": %f, "speed": %f}`, loc.Lat, loc.Lng, loc.Speed)
	
	err := c.client.Set(ctx, redisKey, redisData, 10*time.Minute).Err()
	return err
}

func (c *locationCache) ClearAll(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}
