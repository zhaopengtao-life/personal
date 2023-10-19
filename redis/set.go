package redis

func SetInt(key string, value int64) {
    DbRedis.Set(ctx, key, value, 0)
}
