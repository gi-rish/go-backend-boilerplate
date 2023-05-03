package cache

import (
	"github.com/NewStreetTechnologies/go-backend-boilerplate/config"
	_logger "github.com/NewStreetTechnologies/go-backend-boilerplate/logger"

	_redis "github.com/NewStreetTechnologies/go-backend-common/redis"

	"github.com/go-redis/redis"
)

var (
	logger = _logger.Logger
)

var Cache *redis.Client

func init() {
	logger.Info("Init cache......")
	redis := &_redis.Redis{
		Config: config.GetWholeConfig(),
		Logger: logger,
	}
	host := config.GetConfig("redis.host")
	port := config.GetConfig("redis.port")
	db := config.GetConfig("redis.db")
	password := config.GetConfig("redis.password")
	conn, err := redis.ConnectToRedisClient(host, port, password, db)
	if err != nil {
		logger.Error("Error when connecting to redis")
		panic(err)
	}
	Cache = conn
}
