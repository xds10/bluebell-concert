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
	SeatNo, err := redis.GetSeatByConcertID(p)
	if SeatNo == 0 {
		zap.L().Error("redis.GetSeatByConcertID failed")
		return nil, err
	}
	fmt.Println(SeatNo)
	if p.SeatIdx == nil {
		p.SeatIdx = &models.Seat{} // 假设 SeatIndex 是 SeatIdx 的类型
	}
	p.SeatIdx.SeatID = SeatNo
	err = mysql.BuyTicket(p)
	if err != nil {
		zap.L().Error("mysql.BuyTicket failed", zap.Error(err))
		return nil, err
	}
	// 查询完整的seat表
	p.SeatIdx, err = mysql.GetSeatByID(SeatNo)
	if err != nil {
		zap.L().Error("mysql.GetSeatByID failed", zap.Error(err))
		return nil, err
	}
	order := &models.Order{
		UserID:     p.UserID,
		ConcertID:  p.ConcertID,
		SeatID:     SeatNo,
		Price:      float64(p.SeatIdx.Price),
		Status:     1,
		CreateTime: time.Now(),
	}
	err = mysql.AddOrder(order)
	p.TicketID = order.ID
	if err != nil {
		zap.L().Error("mysql.AddOrders failed", zap.Error(err))
		return nil, err
	}
	return p, nil
}

func PayTicket(p *models.Order) error {
	order := mysql.GetOrderById(p.ID)
	if (time.Now().Sub(order.CreateTime)).Minutes() > 5 {
		zap.L().Error("订单超时")
		return fmt.Errorf("订单超时")
	}
	err := mysql.PayOrder(p.ID)
	if err != nil {
		zap.L().Error("mysql.PayTicket failed", zap.Error(err))
		return err
	}
	return nil
}
