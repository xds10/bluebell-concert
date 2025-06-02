package redis

import (
	"errors"
	"fmt"
	"strconv"

	"bluebell/models"

	"github.com/go-redis/redis"
)

func ConcertExists(key string) (int64, error) {
	return rdb.Exists(GetRedisKey(key)).Result()
}

func AddSeatToConcert(concertID int64, seats *models.Seat) error {
	seatID := strconv.FormatInt(seats.SeatID, 10)
	key := fmt.Sprintf("concert:id:%d:seat:%s", concertID, seats.SeatSection)
	key = GetRedisKey(key)
	_, err := rdb.SAdd(key, seatID).Result()
	if err != nil {
		return fmt.Errorf("添加座位到音乐会失败: %w", err)
	}
	return nil
}

func AddConcertToRedis(concertId int64) error {
	key := fmt.Sprintf("concert:id:%d", concertId)
	key = GetRedisKey(key)
	_, err := rdb.Set(key, concertId, 0).Result()
	if err != nil {
		return fmt.Errorf("添加音乐会到redis失败: %w", err)
	}
	return nil
}

func GetSeatByConcertID(p *models.Ticket) (int64, error) {
	key := fmt.Sprintf("concert:id:%d:seat:%s", p.ConcertID, p.SeatIdx.SeatSection)
	key = GetRedisKey(key)
	seatIDStr, err := rdb.SPop(key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, errors.New("没有可用的座位")
		}
		return 0, err
	}

	// 将座位编号转换为整数
	seatID, err := strconv.ParseInt(seatIDStr, 10, 64)
	if err != nil {
		return 0, errors.New("无效的座位编号")
	}
	return seatID, nil
}
