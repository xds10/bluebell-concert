package logic

import (
	"fmt"

	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"

	"go.uber.org/zap"
)

func CreateConcert(p *models.Concert) (err error) {
	return mysql.InsertConcert(p)
}

func GetConcertList() (data []*models.Concert, err error) {
	return mysql.GetConcertList()
}

func GetConcertDetail(id int64) (data *models.Concert, err error) {
	return mysql.GetConcertByID(id)
}

// GetConcertSeats 获取演唱会座位信息
func GetConcertSeats(concertID int64) (map[string]interface{}, error) {
	// 1. 获取演唱会信息
	concert, err := mysql.GetConcertByID(concertID)
	if err != nil {
		zap.L().Error("mysql.GetConcertByID failed", zap.Error(err))
		return nil, err
	}

	// 2. 获取演唱会座位信息
	seats, err := mysql.GetSeatsByConcertID(concertID)
	if err != nil {
		zap.L().Error("mysql.GetSeatsByConcertID failed", zap.Error(err))
		return nil, err
	}

	// 3. 统计每个区域的总座位数和已选座位数
	sectionStats := make(map[string]map[string]int)
	sectionNames := make([]string, 0)

	for _, seat := range seats {
		section := seat.SeatSection
		if _, ok := sectionStats[section]; !ok {
			sectionStats[section] = map[string]int{
				"total":    0,
				"selected": 0,
			}
			sectionNames = append(sectionNames, section)
		}

		sectionStats[section]["total"]++
		if seat.IsSelected {
			sectionStats[section]["selected"]++
		}
	}

	// 4. 构造返回数据
	sectionsResult := make(map[string]map[string]int)

	for section, stats := range sectionStats {
		total := stats["total"]
		selected := stats["selected"]

		// 直接从Redis获取剩余座位数
		redisKey := fmt.Sprintf("bluebell:concert:id:%d:seat:%s", concertID, section)
		available, err := redis.GetSetSize(redisKey)

		if err != nil {
			zap.L().Warn("Redis获取剩余座位数失败，使用数据库计算",
				zap.Int64("concert_id", concertID),
				zap.String("section", section),
				zap.Error(err))
			// 回退到使用数据库计算的可用座位数
			available = int64(total) - int64(selected)
		} else {
			zap.L().Info("从Redis获取剩余座位数成功",
				zap.Int64("concert_id", concertID),
				zap.String("section", section),
				zap.Int64("available", available))
		}

		sectionsResult[section] = map[string]int{
			"total":     total,
			"available": int(available),
			"selected":  selected,
		}
	}

	result := map[string]interface{}{
		"concert_id":    concertID,
		"concert_title": concert.Title,
		"sections":      sectionsResult,
	}

	return result, nil
}
