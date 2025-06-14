package logic

import (
	"bluebell/dao/mysql"
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

	// 3. 统计每个区域的剩余座位数
	sectionStats := make(map[string]map[string]int)
	for _, seat := range seats {
		section := seat.SeatSection
		if _, ok := sectionStats[section]; !ok {
			sectionStats[section] = map[string]int{
				"total": 0,
				"available": 0,
			}
		}
		
		sectionStats[section]["total"]++
		if !seat.IsSelected {
			sectionStats[section]["available"]++
		}
	}

	// 4. 构造返回数据
	result := map[string]interface{}{
		"concert_id": concertID,
		"concert_title": concert.Title,
		"sections": sectionStats,
	}

	return result, nil
}
