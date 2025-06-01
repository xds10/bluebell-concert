package models

type Seat struct {
	SeatID      int64  `db:"id" json:"id"`
	ConcertId   int64  `db:"concert_id" json:"concert_id"`
	SeatSection string `db:"seat_section" json:"seat_section"`
	SeatRow     int64  `db:"seat_row" json:"seat_row"`
	SeatNo      int64  `db:"seat_No" json:"seat_No"`
	Price       int64  `db:"price" json:"price"`
}
