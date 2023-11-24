package redis

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var (
	ctx             = context.TODO()
	linkRedisMethod sync.Once
	DbRedis         *redis.Client
)

// Init 当程序启动的时候，就初始化链接池
func Init() {
	linkRedisMethod.Do(func() {
		//连接数据库
		DbRedis = redis.NewClient(&redis.Options{
			Addr:     "172.*.*.*:6379", // 对应的ip以及端口号
			Password: "pwdword",        // 数据库的密码
			DB:       6,                // 数据库的编号，默认的话是0
		})
		// 连接测活
		_, err := DbRedis.Ping(ctx).Result()
		if err != nil {
			log.Errorf("InitRedis 连接Redis失败 Error：%v", err)
		}
		log.Info("InitRedis 连接Redis成功")
	})
}
