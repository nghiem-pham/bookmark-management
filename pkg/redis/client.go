package redis

import "github.com/redis/go-redis/v9"

// NewClient creates a redis.Client using configuration loaded from
// environment variables prefixed with envPrefix. It returns an error if
// the configuration cannot be loaded.
func NewClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return redisClient, nil
}
