package redis

import (
	"context"
	"fmt"

	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	redis "github.com/redis/go-redis/v9"
)

type Service interface {
	InitRedis() *redis.Client
	GetClient() *redis.Client
	Ping() error
}

type service struct {
	cfg    config.RedisConfig
	client *redis.Client
}

func NewService(redisConfig config.RedisConfig) Service {
	svc := &service{
		cfg: redisConfig,
	}
	svc.client = svc.InitRedis()
	return svc
}

func (s *service) InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", s.cfg.Hostname, s.cfg.Port),
		Password: s.cfg.Password,
		DB:       s.cfg.DB,
	})

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}

func (s *service) GetClient() *redis.Client {
	return s.client
}

func (s *service) Ping() error {
	return s.client.Ping(context.Background()).Err()
}
