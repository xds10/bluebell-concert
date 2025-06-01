package redis

// redis Key
// redis key 注意使用命名空间区分不同的key,方便查询和拆分
const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"  // zset 帖子及发布时间
	KeyPostScoreZSet   = "post:score" // zset 帖子及投票分数
	KeyPostVotedZSetPF = "post:voted" // zset 记录用户及投票类型 参数是post id
	KeyConcertId       = "concert:id"
)

func GetRedisKey(key string) string {
	return KeyPrefix + key
}
