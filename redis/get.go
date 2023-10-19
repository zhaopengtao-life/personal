package redis

import (
    "github.com/redis/go-redis/v9"
    log "github.com/sirupsen/logrus"
)

func GetInt(key string) int64 {
    offset, err := DbRedis.Get(ctx, key).Int64()
    // redis.Nil 用于判断是否查询到有该组数据
    if err == redis.Nil {
        return offset
    } else if err != nil {
        log.Errorf("GetOffset Redis Get Data Error: %v", err)
        return offset
    } else {
        // 查询成功后处理
        return offset
    }
}
