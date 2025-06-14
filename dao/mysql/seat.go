package mysql

import (
	"bluebell/models"

	"github.com/jmoiron/sqlx"
)

func GetSeatsByVenueID(venueID int) (data []*models.Seat, err error) {
	sqlStr := `select id,section from seat where place_id =?`
	err = db.Select(&data, sqlStr, venueID)
	return
}

func AddSeatToConcert(concertID int64, seatID int64) error {
	sqlStr := `insert into seat_record (concert_id,seat_id) values (?, ?)`
	_, err := db.Exec(sqlStr, concertID, seatID)
	if err != nil {
		return err
	}
	return nil
}

func GetSeatByID(seatID int64) (data *models.Seat, err error) {
	data = &models.Seat{}
	sqlStr := `select id,section,seat_row,seat_no,price from seat where id =?`
	err = db.Get(data, sqlStr, seatID)
	return
}

// GetSeatsByConcertID 获取演唱会的所有座位信息
func GetSeatsByConcertID(concertID int64) ([]*models.SeatRecord, error) {
	sqlStr := `
		SELECT sr.concert_id, sr.seat_id, sr.is_selected, s.section, s.seat_row, s.seat_no, s.price
		FROM seat_record sr
		JOIN seat s ON sr.seat_id = s.id
		WHERE sr.concert_id = ?
	`
	var seats []*models.SeatRecord
	err := db.Select(&seats, sqlStr, concertID)
	if err != nil {
		return nil, err
	}
	return seats, nil
}

// GetSeatByIDTx 在事务中获取座位信息
func GetSeatByIDTx(tx *sqlx.Tx, seatID int64) (data *models.Seat, err error) {
	data = &models.Seat{}
	sqlStr := `select id,section,seat_row,seat_no,price from seat where id =?`
	err = tx.Get(data, sqlStr, seatID)
	return
}

func BuyTicketTx(tx *sqlx.Tx, p *models.Ticket) error {
	// 将座位标记为已选择
	sqlStr := `update seat_record set is_selected = 1 where concert_id = ? and seat_id = ?`
	_, err := tx.Exec(sqlStr, p.ConcertID, p.SeatIdx.SeatID)
	return err
}

// GetSeatSectionBySeatID 获取座位的区域信息
func GetSeatSectionBySeatID(seatID int64) (string, error) {
	sqlStr := `select section from seat where id = ?`
	var section string
	err := db.Get(&section, sqlStr, seatID)
	return section, err
}
