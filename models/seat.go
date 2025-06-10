package models

type Seat struct {
	SeatID      int64   `db:"id" json:"id"`
	ConcertId   int64   `db:"concert_id" json:"concert_id"`
	SeatSection string  `db:"section" json:"section"`
	SeatRow     int64   `db:"seat_row" json:"seat_row"`
	SeatNo      int64   `db:"seat_no" json:"seat_no"`
	Price       float64 `db:"price" json:"price"`
	Place       string  `db:"place" json:"place"`
	PlaceID     int64   `db:"place_id" json:"place_id"`
}
