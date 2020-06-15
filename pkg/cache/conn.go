package cache

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/sunnywalden/sync-data/pkg/logging"

	"github.com/go-redis/redis/v8"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/errors"
)

var (
	log *logrus.Logger
)

// GetClient, init redis client
func GetClient(ctx context.Context,configures *config.TomlConfig) ( *redis.Client, error) {

	log = logging.GetLogger(&configures.Log)

	redisConf := configures.Redis
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

// Close, mysql client closing
func Close(client * redis.Client) error {
	err := client.Close()
	return err
}


