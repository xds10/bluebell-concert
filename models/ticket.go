package models

type Ticket struct {
	TicketID  int64 `db:"id" json:"id"`
	ConcertID int64 `db:"concert_id" json:"concert_id"`
	UserID    int64 `db:"user_id" json:"user_id"`
	SeatIdx   *Seat `db:"seat_idx" json:"seat_idx"`
}
