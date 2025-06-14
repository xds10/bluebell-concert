package redis

import (
	"errors"
	"fmt"
	"strconv"

	"bluebell/models"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
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

// ReturnSeatToConcert 将座位放回演唱会的可用座位集合
func ReturnSeatToConcert(concertID int64, seatID int64, section string) error {
	key := fmt.Sprintf("concert:id:%d:seat:%s", concertID, section)
	key = GetRedisKey(key)
	seatIDStr := strconv.FormatInt(seatID, 10)
	
	// 将座位ID添加回集合
	_, err := rdb.SAdd(key, seatIDStr).Result()
	if err != nil {
		return fmt.Errorf("将座位放回集合失败: %w", err)
	}
	
	return nil
}

// GetRemainingSeats 获取演唱会某个区域的剩余座位数
func GetRemainingSeats(concertID int64, section string) (int64, error) {
	key := fmt.Sprintf("concert:id:%d:seat:%s", concertID, section)
	key = GetRedisKey(key)
	
	// 使用SCARD命令获取集合中元素的数量
	count, err := rdb.SCard(key).Result()
	if err != nil {
		return 0, fmt.Errorf("获取剩余座位数失败: %w", err)
	}
	
	return count, nil
}

// GetAllSectionsRemainingSeats 获取演唱会所有区域的剩余座位数
func GetAllSectionsRemainingSeats(concertID int64, sections []string) (map[string]int64, error) {
	result := make(map[string]int64)
	
	for _, section := range sections {
		count, err := GetRemainingSeats(concertID, section)
		if err != nil {
			// 记录错误但继续处理其他区域
			zap.L().Error("获取区域剩余座位数失败", 
				zap.Int64("concert_id", concertID),
				zap.String("section", section),
				zap.Error(err))
			result[section] = 0
		} else {
			result[section] = count
		}
	}
	
	return result, nil
}

// GetSectionTotalSeats 获取演唱会某个区域的总座位数
// 这个函数应该在初始化时将总座位数存入Redis
func GetSectionTotalSeats(concertID int64, section string) (int64, error) {
	key := fmt.Sprintf("concert:id:%d:section:%s:total", concertID, section)
	key = GetRedisKey(key)
	
	// 从Redis获取总座位数
	count, err := rdb.Get(key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("演唱会区域总座位数未初始化")
		}
		return 0, fmt.Errorf("获取总座位数失败: %w", err)
	}
	
	return count, nil
}

// SetSectionTotalSeats 设置演唱会某个区域的总座位数
func SetSectionTotalSeats(concertID int64, section string, total int64) error {
	key := fmt.Sprintf("concert:id:%d:section:%s:total", concertID, section)
	key = GetRedisKey(key)
	
	// 将总座位数存入Redis
	_, err := rdb.Set(key, total, 0).Result()
	if err != nil {
		return fmt.Errorf("设置总座位数失败: %w", err)
	}
	
	return nil
}

// GetConcertSections 获取演唱会的所有区域
func GetConcertSections(concertID int64) ([]string, error) {
	key := fmt.Sprintf("concert:id:%d:sections", concertID)
	key = GetRedisKey(key)
	
	// 从Redis获取所有区域
	sections, err := rdb.SMembers(key).Result()
	if err != nil {
		return nil, fmt.Errorf("获取演唱会区域失败: %w", err)
	}
	
	return sections, nil
}

// AddConcertSection 添加演唱会区域
func AddConcertSection(concertID int64, section string) error {
	key := fmt.Sprintf("concert:id:%d:sections", concertID)
	key = GetRedisKey(key)
	
	// 将区域添加到Redis集合
	_, err := rdb.SAdd(key, section).Result()
	if err != nil {
		return fmt.Errorf("添加演唱会区域失败: %w", err)
	}
	
	return nil
}

// GetConcertSeatInfo 只从Redis获取演唱会座位信息
func GetConcertSeatInfo(concertID int64) (map[string]map[string]int64, error) {
	// 1. 获取演唱会的所有区域
	sections, err := GetConcertSections(concertID)
	if err != nil {
		return nil, fmt.Errorf("获取演唱会区域失败: %w", err)
	}
	
	// 如果没有区域信息，返回空结果
	if len(sections) == 0 {
		return make(map[string]map[string]int64), nil
	}
	
	// 2. 获取每个区域的信息
	result := make(map[string]map[string]int64)
	
	for _, section := range sections {
		// 获取区域的总座位数
		total, err := GetSectionTotalSeats(concertID, section)
		if err != nil {
			zap.L().Warn("获取区域总座位数失败",
				zap.Int64("concert_id", concertID),
				zap.String("section", section),
				zap.Error(err))
			total = 0
		}
		
		// 获取区域的剩余座位数
		available, err := GetRemainingSeats(concertID, section)
		if err != nil {
			zap.L().Warn("获取区域剩余座位数失败",
				zap.Int64("concert_id", concertID),
				zap.String("section", section),
				zap.Error(err))
			available = 0
		}
		
		// 计算已选座位数
		selected := total - available
		if selected < 0 {
			selected = 0 // 确保不会出现负数
		}
		
		// 将结果添加到返回值
		result[section] = map[string]int64{
			"total": total,
			"available": available,
			"selected": selected,
		}
	}
	
	return result, nil
}

// InitConcertSeatInfo 初始化演唱会座位信息到Redis
func InitConcertSeatInfo(concertID int64, sectionStats map[string]int64) error {
	// 添加演唱会的所有区域
	for section, total := range sectionStats {
		// 添加区域到演唱会区域集合
		err := AddConcertSection(concertID, section)
		if err != nil {
			return fmt.Errorf("添加演唱会区域失败: %w", err)
		}
		
		// 设置区域的总座位数
		err = SetSectionTotalSeats(concertID, section, total)
		if err != nil {
			return fmt.Errorf("设置区域总座位数失败: %w", err)
		}
	}
	
	return nil
}

// GetSetSize 直接获取指定Redis集合的大小
func GetSetSize(key string) (int64, error) {
	// 使用SCARD命令获取集合中元素的数量
	count, err := rdb.SCard(key).Result()
	if err != nil {
		return 0, fmt.Errorf("获取集合大小失败: %w", err)
	}
	
	return count, nil
}
