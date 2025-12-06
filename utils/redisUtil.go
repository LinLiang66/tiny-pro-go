package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisUtil struct {
	client *redis.Client
}

// Redis  全局变量, 外部使用utils.Redis来访问
var Redis RedisUtil

// 初始化redis
func init() {
	viper.SetConfigFile("./config/config.yaml")

	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		panic(err)
		return
	}
	//连接redis
	r := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", viper.GetString("cache.redis.host"), viper.GetInt("cache.redis.port")),
		Password:    viper.GetString("cache.redis.password"),
		DB:          viper.GetInt("cache.redis.db"),
		PoolSize:    viper.GetInt("cache.redis.pool_size"),
		ReadTimeout: time.Duration(viper.GetInt("cache.redis.timeout")),
	})
	//初始化全局redis连接
	Redis = RedisUtil{client: r}
}

// SetStr 设置数据到redis中（string）
func (rs *RedisUtil) SetStr(ctx context.Context, key string, value string, expiration time.Duration) error {
	_, err := rs.client.Set(ctx, key, value, expiration).Result()
	return err
}

// SetStrNotExist 设置数据到redis中（string）
func (rs *RedisUtil) SetStrNotExist(ctx context.Context, key string, value string, expireSecond int) bool {
	val, err := rs.client.Do(ctx, "SET", key, value, "EX", expireSecond, "NX").Result()
	if err != nil || val == nil {
		return false
	}
	return true
}

// SetEx 设置数据到redis中
func (rs *RedisUtil) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rs.client.SetEx(ctx, key, value, expiration).Err()
}

// GetStr 获取redis中数据（string）
func (rs *RedisUtil) GetStr(ctx context.Context, key string) (string, error) {
	val, err := rs.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// HSet 设置数据到redis中（hash）
func (rs *RedisUtil) HSet(ctx context.Context, key string, field string, value string) error {
	return rs.client.Do(ctx, "HSet", key, field, value).Err()
}

// HGet 获取redis中数据（hash）
func (rs *RedisUtil) HGet(ctx context.Context, key string, field string) (string, error) {
	val, err := rs.client.Do(ctx, "HGet", key, field).Result()
	if err != nil {
		return "", err
	}
	return string(val.([]byte)), nil
}

// DelByKey 删除
func (rs *RedisUtil) DelByKey(ctx context.Context, key string) error {
	return rs.client.Del(ctx, "DEL", key).Err()

}

// SetExpire 设置key过期时间
func (rs *RedisUtil) SetExpire(ctx context.Context, key string, expiration time.Duration) error {
	return rs.client.Do(ctx, "EXPIRE", key, expiration).Err()
}

// Exists 判断KEY在redis中是否存在
func (rs *RedisUtil) Exists(ctx context.Context, KEY string) bool {
	exists, err := rs.client.Do(ctx, "EXISTS", KEY).Bool()
	if err != nil {
		return false
	}
	return exists
}

// KEYEXISTSGetStr 判断KEY在redis中是否存在,存在则获取内容
func (rs *RedisUtil) KEYEXISTSGetStr(ctx context.Context, KEY string) (bool, string) {
	if rs.Exists(ctx, KEY) {
		str, err := rs.GetStr(ctx, KEY)
		if err == nil {
			return true, str
		}
	}
	return false, ""
}

// GetBytes  获取redis中数据（string）
func (rs *RedisUtil) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return rs.client.Get(ctx, key).Bytes()
}

// KEYEXISTSGetBytes  判断KEY在redis中是否存在,存在则获取内容
func (rs *RedisUtil) KEYEXISTSGetBytes(ctx context.Context, KEY string) (bool, []byte) {
	if rs.Exists(ctx, KEY) {
		str, err := rs.GetBytes(ctx, KEY)
		if err == nil {
			return true, str
		}
	}
	return false, nil
}

// KEYEXISTSGetScan  判断KEY在redis中是否存在,存在则获取指定类型的内容
func (rs *RedisUtil) KEYEXISTSGetScan(ctx context.Context, KEY string, Val interface{}) bool {
	if rs.Exists(ctx, KEY) {
		body, err := rs.client.Get(ctx, KEY).Bytes()
		if err == nil {
			if json.Unmarshal(body, Val) == nil {
				return true
			}
		}
		log.Println("获取redis 数据报错了", err.Error())
	}
	return false
}

// Set 设置数据到redis中 泛型
func (rs *RedisUtil) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rs.client.Set(ctx, key, value, expiration).Err()
}

// Keys  获取Redis中特定路径前缀的所有键
func (rs *RedisUtil) Keys(ctx context.Context, key string) ([]string, error) {
	return rs.client.Keys(ctx, key).Result()
}

// LeftPushJSON 将给定的 JSON 字符串插入到 Redis 列表的左侧
func (rs *RedisUtil) LeftPushJSON(key string, value interface{}) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Printf("Failed to marshal value to JSON: %v", err)
		return
	}
	err = rs.client.LPush(context.Background(), key, jsonValue).Err()
	if err != nil {
		log.Printf("Failed to left push JSON value to list: %v", err)
		return
	}

}
