package logic

import (
	"fmt"

	"bluebell/dao/mysql"

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

		zap.L().Info("Order expired automatically",
			zap.Int64("order_id", order.ID),
			zap.Time("create_time", order.CreateTime))
	}
}
