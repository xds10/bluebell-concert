package models

type Concert struct {
	ID       int64  `json:"id" db:"id"`
	Title    string `json:"title" db:"title"`
	VenueId  int    `json:"venue_id" db:"venue_id"`
	Date     string `json:"date" db:"start_time"`
	EndDate  string `json:"enddate" db:"end_time"`
	SaleTime string `json:"saletime" db:"release_time"`
	Status   int    `json:"status" db:"status"`
}
