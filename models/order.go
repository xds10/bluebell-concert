package models

import "time"

type Order struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	ConcertID  int64     `json:"concert_id" db:"concert_id"`
	SeatID     int64     `json:"seat_id" db:"seat_id"`
	Price      float64   `json:"price" db:"price"`
	Status     int       `json:"status" db:"status"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	Ticket     *Ticket   `json:"ticket" db:"ticket"`
}
