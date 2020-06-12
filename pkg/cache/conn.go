package cache

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/errors"
	"github.com/sunnywalden/sync-data/pkg/logging"
)

var (
	log = logging.GetLogger()
)

// GetClient, init redis client
func GetClient(ctx context.Context,configures *config.RedisConf) ( *redis.Client, error) {

	redisConf := configures
	redisHost := redisConf.Host + ":" + redisConf.Port
	redisDB   := redisConf.DB
	redisPassword := redisConf.Password

	log.Debugf("Debug redis addr:%s\n", redisHost)

	options := &redis.Options{
		Addr:     redisHost,
		DB:       redisDB,
		Password: redisPassword,
	}

	if options.Addr == "" {
		return nil, errors.ErrRedisConfigNil
	}

	client := redis.NewClient(options)
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	} else {
		log.Infof("redis ping response:%s\n", pong)
	}

	return client, nil
}

func Close(client * redis.Client) error {
	err := client.Close()
	return err
}


