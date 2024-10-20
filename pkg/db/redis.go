package db

import (
	"IM/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	userOnlinePrefix = "user_online_" // 用户在线状态设置
	ttl1D            = 24 * 60 * 60   // s  1天

	seqPrefix         = "object_seq_" // 群成员信息
	SeqObjectTypeUser = 1             // 用户

)

var (
	RDB *redis.Client
)

// Redis初始化
func InitRedis(cfg *config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		DB:           cfg.DB,
		Password:     cfg.Password,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	return nil
}

func getUserKey(userId int64) string {
	return fmt.Sprintf("%s%d", userOnlinePrefix, userId)
}

// SetUserOnline 设置用户在线
func SetUserOnline(userId int64, addr string) error {
	key := getUserKey(userId)
	_, err := RDB.Set(context.Background(), key, addr, ttl1D*time.Second).Result()
	if err != nil {
		fmt.Println("[设置用户在线] 错误, err:", err)
		return err
	}
	return nil
}

// GetUserOnline 获取userId的在线状态
func GetUserOnline(userId int64) (string, error) {
	key := getUserKey(userId)
	addr, err := RDB.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("[获取用户在线] 错误，err:", err)
		return "", err
	}
	return addr, nil
}

// DelUserOnline 删除用户在线信息（存在即在线）
func DelUserOnline(userId int64) error {
	key := getUserKey(userId)
	_, err := RDB.Del(context.Background(), key).Result()
	if err != nil {
		fmt.Println("[删除用户在线] 错误, err:", err)
		return err
	}
	return nil
}

func getSeqKey(objectType int8, userID int64) string {
	return fmt.Sprintf("%s%d_%d", seqPrefix, objectType, userID)
}

// 获取用户的下一个seq，消息同步序列号
func GetUserNextSeq(objectType int8, objectID int64) (int64, error) {
	key := getSeqKey(objectType, objectID)
	result, err := RDB.Incr(context.Background(), key).Uint64()
	if err != nil {
		fmt.Println("GetUserNextSeq.RDB.Incr error:", err)
		return 0, err
	}
	return int64(result), nil
}

func GetUserNextSeqBatch(objectType int8, objectIDs []int64) ([]int64, error) {
	// 在从Redis获取一组数据的时候，不要单一的使用循环，这会大大增加了RTT时间
	script := `
       local results = {}
       for i, key in ipairs(KEYS) do
           results[i] = redis.call('INCR', key)
       end
       return results
   `
	keys := make([]string, len(objectIDs))
	for i, objectID := range objectIDs {
		keys[i] = getSeqKey(objectType, objectID)
	}
	// 使用Lua脚本确保操作的原子性、减少网络通信、简化逻辑以及提升性能
	res, err := redis.NewScript(script).Run(context.Background(), RDB, keys).Result()
	if err != nil {
		fmt.Println("[获取seq] 失败，err:", err)
		return nil, err
	}
	results := make([]int64, len(objectIDs))
	for i, v := range res.([]interface{}) {
		results[i] = v.(int64)
	}
	return results, nil
}
