package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"

	"go.uber.org/zap"
)

func BuyTicket(p *models.Ticket) error {
	SeatNo, err := redis.GetSeatByConcertID(p.ConcertID)
	if SeatNo == 0 {
		zap.L().Error("redis.GetSeatByConcertID failed")
		return err
	}
	p.SeatIdx.SeatID = SeatNo
	err = mysql.BuyTicket(p)
	return err
}
