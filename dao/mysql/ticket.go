package mysql

import (
	"bluebell/models"
)

func BuyTicket(p *models.Ticket) error {
	sqlStr := `update seat_record set is_selected = 1 where concert_id=? and seat_id=?`
	_, err := db.Exec(sqlStr, p.ConcertID, p.SeatIdx.SeatID)
	return err
}
func GetTicketByOrderID(orderID int64) (*models.Ticket, error) {
	sqlStr := `select id,place_id,place,section,seat_row,seat_no from seat where id=?`
	ticket := &models.Ticket{}
	ticket.SeatIdx = &models.Seat{}
	err := db.Get(ticket.SeatIdx, sqlStr, orderID)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}