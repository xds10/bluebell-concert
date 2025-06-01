package mysql

import "bluebell/models"

func InsertConcert(p *models.Concert) (err error) {
	sqlStr := `insert into concert( title, venue_id, start_time,end_time, release_time, status) values (?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.Title, p.VenueId, p.Date, p.EndDate, p.SaleTime, p.Status)
	return err
}

func GetConcertList() (data []*models.Concert, err error) {
	sqlStr := `select id,title, venue_id, start_time, end_time, release_time, status from concert`
	err = db.Select(&data, sqlStr)
	return
}

func GetConcertByID(id int64) (data *models.Concert, err error) {
	data = new(models.Concert)
	sqlStr := `select id,title, venue_id, start_time, end_time, release_time, status from concert where id = ?`
	err = db.Get(data, sqlStr, id)
	return data, err
}

func GetFutureConcerts() (data []*models.Concert, err error) {
	sqlStr := `select id,title, venue_id, start_time, end_time, release_time, status from concert where start_time > now()`
	err = db.Select(&data, sqlStr)
	return
}
