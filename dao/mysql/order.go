package mysql

import "bluebell/models"

func AddOrder(order *models.Order) error {
	sqlStr := `insert into orders (user_id,concert_id,seat_id,price,status,create_time) values (?,?,?,?,?,?)`
	result, err := db.Exec(sqlStr, order.UserID, order.ConcertID, order.SeatID, order.Price, order.Status, order.CreateTime)
	if err != nil {
		return err
	}
	order.ID, err = result.LastInsertId()
	return err
}

func PayOrder(id int64) error {
	sqlStr := `update orders set status = ? where id = ?`
	_, err := db.Exec(sqlStr, 2, id)
	return err
}

func GetOrderById(id int64) *models.Order {
	sqlStr := `select id,user_id,concert_id,seat_id,price,status,create_time from orders where id =?`
	order := &models.Order{}
	err := db.Get(order, sqlStr, id)
	if err != nil {
		return nil
	}
	return order
}

func GetOrderList(userID int64) ([]*models.Order, error) {
	sqlStr := `select id,user_id,concert_id,seat_id,price,status,create_time from orders where user_id = ?`
	orders := []*models.Order{}
	err := db.Select(&orders, sqlStr, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}