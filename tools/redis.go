package tools

import (
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"time"
)

const (
	Nil = redis.Nil
)

var Redis = &redisCli{}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type redisCli struct {
	*redis.Client
}

// InitRedis @Title 初始化redis
func InitRedis(config *RedisConfig) (string, error) {
	Redis.Client = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // no password set
		DB:       config.DB,       // use default DB
	})
	return Redis.Ping().Result()
}

// SetStruct @Title 序列化
func (r *redisCli) SetStruct(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	marshal, _ := msgpack.Marshal(value)
	return r.Set(key, marshal, expiration)
}

// GetStruct @Title 反序列化
func (r *redisCli) GetStruct(key string, val interface{}) error {
	return msgpack.Unmarshal([]byte(r.Get(key).Val()), val)
}
