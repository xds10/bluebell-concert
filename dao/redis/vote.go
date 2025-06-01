package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"math"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 60 * 60
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("vote time expire")
)

func CreatePost(postID int64) error {

	pipeline := rdb.Pipeline() //事务
	key := GetRedisKey(KeyPostTimeZSet)

	// 检查键的类型
	keyType, err := rdb.Type(key).Result()
	if err != nil {
		zap.L().Error("check key type failed", zap.Error(err))
		return err
	}

	// 如果键的类型不是有序集合，删除该键
	if keyType != "zset" {
		err := rdb.Del(key).Err()
		if err != nil {
			zap.L().Error("delete key failed", zap.Error(err))
			return err
		}
	}

	// 添加帖子 ID 到有序集合
	_, err = pipeline.ZAdd(key, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	}).Result()
	if err != nil {
		zap.L().Error("create post failed", zap.Error(err))
	}
	_, err = pipeline.ZAdd(GetRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	}).Result()

	_, err = pipeline.Exec()
	return err
}
func VoteForPost(userID, postID string, direction float64) error {
	// 1.判断投票限制
	// 去redis取帖子发布时间
	postTime := rdb.ZScore(GetRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		zap.L().Debug("time is over",
			zap.Float64("postTime", postTime),
			zap.Float64("time", float64(time.Now().Unix())))
		return ErrVoteTimeExpire
	}
	// 2.更新分数
	// 先查之前的投票记录
	ov := rdb.ZScore(GetRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	diff := math.Abs((ov - direction))
	var dir float64
	if direction > ov {
		dir = 1
	} else {
		dir = -1
	}
	_, err := rdb.ZIncrBy(GetRedisKey(KeyPostScoreZSet), dir*diff*scorePerVote, postID).Result()
	if err != nil {
		return err
	}
	// 3.记录用户为该帖子投票
	if direction == 0 {
		rdb.ZRem(GetRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		rdb.ZAdd(GetRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  direction,
			Member: userID,
		}).Result()
	}

	return nil
}
