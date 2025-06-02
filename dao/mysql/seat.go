package mysql

import "bluebell/models"

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
