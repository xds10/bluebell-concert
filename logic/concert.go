package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func CreateConcert(p *models.Concert) (err error) {
	return mysql.InsertConcert(p)
}

func GetConcertList() (data []*models.Concert, err error) {
	return mysql.GetConcertList()
}

func GetConcertDetail(id int64) (data *models.Concert, err error) {
	return mysql.GetConcertByID(id)
}
