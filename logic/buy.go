package logic

import (
	"fmt"
	"time"

	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"

	"go.uber.org/zap"
)

func BuyTicket(p *models.Ticket) (*models.Ticket, error) {
	// 1. 从Redis获取座位
	SeatNo, err := redis.GetSeatByConcertID(p)
	if SeatNo == 0 {
		zap.L().Error("redis.GetSeatByConcertID failed", zap.Error(err))
		return nil, fmt.Errorf("获取座位失败: %w", err)
	}
	zap.L().Info("获取到座位", zap.Int64("座位ID", SeatNo))
	
	if p.SeatIdx == nil {
		p.SeatIdx = &models.Seat{}
	}
	p.SeatIdx.SeatID = SeatNo
	
	// 2. 开始数据库事务
	tx, err := mysql.DB().Beginx()
	if err != nil {
		zap.L().Error("begin transaction failed", zap.Error(err))
		return nil, fmt.Errorf("开始数据库事务失败: %w", err)
	}
	
	// 声明事务结果变量
	var txErr error
	
	// 使用defer处理事务提交或回滚
	defer func() {
		if p := recover(); p != nil {
			// 发生恐慌时回滚事务
			_ = tx.Rollback()
			panic(p) // 继续向上传递恐慌
		} else if txErr != nil {
			// 发生错误时回滚事务
			_ = tx.Rollback()
		}
	}()
	
	// 3. 在事务中执行购票操作
	if txErr = mysql.BuyTicketTx(tx, p); txErr != nil {
		zap.L().Error("mysql.BuyTicket failed", zap.Error(txErr))
		return nil, fmt.Errorf("购票操作失败: %w", txErr)
	}
	
	// 4. 在事务中查询完整的seat表
	p.SeatIdx, txErr = mysql.GetSeatByIDTx(tx, SeatNo)
	if txErr != nil {
		zap.L().Error("mysql.GetSeatByID failed", zap.Error(txErr))
		return nil, fmt.Errorf("获取座位信息失败: %w", txErr)
	}
	
	// 5. 在事务中创建订单
	order := &models.Order{
		UserID:     p.UserID,
		ConcertID:  p.ConcertID,
		SeatID:     SeatNo,
		Price:      float64(p.SeatIdx.Price),
		Status:     1,
		CreateTime: time.Now(),
	}
	
	txErr = mysql.AddOrderTx(tx, order)
	if txErr != nil {
		zap.L().Error("mysql.AddOrders failed", zap.Error(txErr))
		return nil, fmt.Errorf("创建订单失败: %w", txErr)
	}
	
	// 提交事务
	if txErr = tx.Commit(); txErr != nil {
		zap.L().Error("commit transaction failed", zap.Error(txErr))
		return nil, fmt.Errorf("提交事务失败: %w", txErr)
	}
	
	p.TicketID = order.ID
	return p, nil
}

func PayTicket(p *models.Order) error {
	order, err := mysql.GetOrderById(p.ID)
	if (time.Now().Sub(order.CreateTime)).Minutes() > 5 {
		zap.L().Error("订单超时")
		return fmt.Errorf("订单超时")
	}
	err = mysql.PayOrder(p.ID)
	if err != nil {
		zap.L().Error("mysql.PayTicket failed", zap.Error(err))
		return err
	}
	return nil
}

func OrderList(user *models.User) (data []*models.Order, err error) {
	orders, err := mysql.GetOrderList(user.UserID)
	if err != nil {
		zap.L().Error("mysql.GetOrderList failed", zap.Error(err))
		return nil, err
	}
	return orders, nil
}

func GetOrderDetail(p int) (*models.Order, error) {
	order, err := mysql.GetOrderById(int64(p))
	if order == nil {
		return nil, fmt.Errorf("订单不存在")
	}
	order.Ticket = &models.Ticket{
		TicketID:  order.ID,
		ConcertID: order.ConcertID,
		UserID:    order.UserID,
	}
	order.Ticket, err = mysql.GetTicketByOrderID(order.SeatID)
	if err != nil {
		zap.L().Error("mysql.GetTicketByOrderID failed", zap.Error(err))
		return nil, err
	}
	return order, nil
}
