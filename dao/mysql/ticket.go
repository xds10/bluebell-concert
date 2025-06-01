package mysql

import (
	"bluebell/models"
)

func BuyTicket(p *models.Ticket) error {
	sqlStr := `update seat_record set is_selected = 1 where concert_id=? and seat_id=?`
	_, err := db.Exec(sqlStr, p.ConcertID, p.SeatIdx)
	return err
}
