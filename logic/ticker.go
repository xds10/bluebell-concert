package logic

import (
	"fmt"
	"time"

	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"

	"go.uber.org/zap"
)

func Ticker(control chan string, done chan bool) {
	// 创建两个定时器，一个用于演唱会排班检查（2小时），一个用于订单过期检查（1分钟）
	concertTicker := time.NewTicker(2 * time.Hour)
	orderTicker := time.NewTicker(1 * time.Minute)
	
	defer concertTicker.Stop()
	defer orderTicker.Stop()
	
	// 初始执行一次任务
	executeConcertTask()
	executeOrderExpirationTask()

	for {
		select {
		case <-concertTicker.C:
			// 每2小时触发一次演唱会排班任务
			executeConcertTask()
		case <-orderTicker.C:
			// 每1分钟触发一次订单过期检查
			executeOrderExpirationTask()
		case cmd := <-control:
			// 处理控制命令
			switch cmd {
			case "stop":
				fmt.Println("收到停止命令，定时器即将停止")
				done <- true
				return
			case "reset":
				fmt.Println("收到重置命令，重新启动定时器")
				concertTicker.Stop()
				orderTicker.Stop()
				concertTicker = time.NewTicker(2 * time.Hour)
				orderTicker = time.NewTicker(1 * time.Minute)
			}
		}
	}
}

// 执行订单过期检查任务
func executeOrderExpirationTask() {
	zap.L().Info("执行订单过期检查任务...")
	CheckExpiredOrders()
}

// 执行演唱会排班任务
func executeConcertTask() {
	zap.L().Info("执行演唱会排班任务...")
	// 读取未来1个月的演唱会排班
	concerts, err := mysql.GetFutureConcerts()
	if err != nil {
		zap.L().Error("获取未来1个月的演唱会排班失败", zap.Error(err))
		return
	}
	// 遍历每个演唱会
	futureConcert := make([]*models.Concert, 0)
	for _, concert := range concerts {
		// 检查redis中是否存在该演唱会的座位信息
		exists, err := redis.ConcertExists(fmt.Sprintf("concert:id:%d", concert.ID))
		if err != nil {
			zap.L().Error("检查redis中是否存在该演唱会的座位信息失败", zap.Error(err))
			return
		}
		// fmt.Println(exists)
		if exists != 0 {
			// 如果存在，跳过该演唱会
			continue
		}
		futureConcert = append(futureConcert, concert)
	}
	// 找到还没在redis的concert:id记录的演唱会
	// 将演唱会对应的场馆的座位导入到redis的set中
	for _, concert := range futureConcert {

		// 获取场馆的座位信息
		seats, err := mysql.GetSeatsByVenueID(concert.VenueId)
		if err != nil {
			zap.L().Error("获取场馆的座位信息失败", zap.Error(err))
			return
		}
		// 将座位信息导入到redis的set中
		for _, seat := range seats {
			err = redis.AddSeatToConcert(concert.ID, seat)
			if err != nil {
				zap.L().Error("将座位信息导入到redis的set中失败", zap.Error(err))
				return
			}
			err = mysql.AddSeatToConcert(concert.ID, seat.SeatID)
			if err != nil {
				zap.L().Error("将座位信息导入到mysql的set中失败", zap.Error(err))
				return
			}
		}
		err = redis.AddConcertToRedis(concert.ID)
		if err != nil {
			zap.L().Error("将座位信息导入到redis的set中失败", zap.Error(err))
			return
		}
	}
}

