package logic

import (
	"fmt"

	"bluebell/dao/mysql"
	"bluebell/dao/redis"

	"go.uber.org/zap"
)

func CancelOrder(orderID int64) error {
	order, err := mysql.GetOrderById(orderID)
	if err != nil {
		zap.L().Error("mysql.GetOrderById failed", zap.Error(err))
		return err
	}

	// 只有未支付的订单才能取消
	if order.Status != 1 {
		return fmt.Errorf("只有未支付的订单可以取消")
	}

	err = mysql.CancelOrder(orderID)
	if err != nil {
		zap.L().Error("mysql.CancelOrder failed", zap.Error(err))
		return err
	}

	// 释放座位
	err = mysql.ReleaseSeat(order.ConcertID, order.SeatID)
	if err != nil {
		zap.L().Error("mysql.ReleaseSeat failed", zap.Error(err))
		// 即使释放座位失败，也不影响订单取消
	}

	// 获取座位区域信息
	section, err := mysql.GetSeatSectionBySeatID(order.SeatID)
	if err != nil {
		zap.L().Error("mysql.GetSeatSectionBySeatID failed", zap.Error(err))
		// 继续处理，不影响订单取消
	} else {
		// 将座位放回Redis集合
		err = redis.ReturnSeatToConcert(order.ConcertID, order.SeatID, section)
		if err != nil {
			zap.L().Error("redis.ReturnSeatToConcert failed", zap.Error(err))
			// 继续处理，不影响订单取消
		} else {
			zap.L().Info("座位已成功放回Redis集合",
				zap.Int64("concert_id", order.ConcertID),
				zap.Int64("seat_id", order.SeatID),
				zap.String("section", section))
		}
	}

	return nil
}

// CheckExpiredOrders 检查并处理过期订单
func CheckExpiredOrders() {
	// 获取所有过期未支付订单
	expiredOrders, err := mysql.GetExpiredOrders()
	if err != nil {
		zap.L().Error("mysql.GetExpiredOrders failed", zap.Error(err))
		return
	}

	// 处理每个过期订单
	for _, order := range expiredOrders {
		// 更新订单状态为已过期
		err := mysql.ExpireOrder(order.ID)
		if err != nil {
			zap.L().Error("mysql.ExpireOrder failed",
				zap.Int64("order_id", order.ID),
				zap.Error(err))
			continue
		}

		// 释放座位
		err = mysql.ReleaseSeat(order.ConcertID, order.SeatID)
		if err != nil {
			zap.L().Error("mysql.ReleaseSeat failed",
				zap.Int64("concert_id", order.ConcertID),
				zap.Int64("seat_id", order.SeatID),
				zap.Error(err))
			// 即使释放座位失败，也不影响订单过期
		}

		// 获取座位区域信息
		section, err := mysql.GetSeatSectionBySeatID(order.SeatID)
		if err != nil {
			zap.L().Error("mysql.GetSeatSectionBySeatID failed", 
				zap.Int64("seat_id", order.SeatID),
				zap.Error(err))
		} else {
			// 将座位放回Redis集合
			err = redis.ReturnSeatToConcert(order.ConcertID, order.SeatID, section)
			if err != nil {
				zap.L().Error("redis.ReturnSeatToConcert failed", 
					zap.Int64("concert_id", order.ConcertID),
					zap.Int64("seat_id", order.SeatID),
					zap.Error(err))
			} else {
				zap.L().Info("过期订单座位已成功放回Redis集合",
					zap.Int64("concert_id", order.ConcertID),
					zap.Int64("seat_id", order.SeatID),
					zap.String("section", section))
			}
		}

		zap.L().Info("Order expired automatically",
			zap.Int64("order_id", order.ID),
			zap.Time("create_time", order.CreateTime))
	}
}
