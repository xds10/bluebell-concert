package mysql

import (
	"time"
	
	"bluebell/models"
	
	"github.com/jmoiron/sqlx"
)

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

func CancelOrder(id int64) error {
	sqlStr := `update orders set status = ? where id = ?`
	_, err := db.Exec(sqlStr, 4, id)
	return err
}

func ExpireOrder(id int64) error {
	sqlStr := `update orders set status = ? where id = ?`
	_, err := db.Exec(sqlStr, 3, id)
	return err
}

func GetOrderById(id int64) (*models.Order, error) {
	sqlStr := `select id,user_id,concert_id,seat_id,price,status,create_time from orders where id =?`
	order := &models.Order{}
	err := db.Get(order, sqlStr, id)
	if err != nil {
		return nil, err
	}
	return order, nil
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

// 获取所有未支付且创建时间超过5分钟的订单
func GetExpiredOrders() ([]*models.Order, error) {
	sqlStr := `select id,user_id,concert_id,seat_id,price,status,create_time from orders 
               where status = 1 and create_time < ?`
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	orders := []*models.Order{}
	err := db.Select(&orders, sqlStr, fiveMinutesAgo)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// AddOrderTx 在事务中添加订单
func AddOrderTx(tx *sqlx.Tx, order *models.Order) error {
	sqlStr := `insert into orders (user_id,concert_id,seat_id,price,status,create_time) values (?,?,?,?,?,?)`
	result, err := tx.Exec(sqlStr, order.UserID, order.ConcertID, order.SeatID, order.Price, order.Status, order.CreateTime)
	if err != nil {
		return err
	}
	order.ID, err = result.LastInsertId()
	return err
}

// GetOrderBySeat 根据演唱会ID和座位ID查询订单
func GetOrderBySeat(concertID, seatID int64) (*models.Order, error) {
	sqlStr := `select id,user_id,concert_id,seat_id,price,status,create_time 
              from orders where concert_id = ? and seat_id = ? 
              order by create_time desc limit 1`
	order := &models.Order{}
	err := db.Get(order, sqlStr, concertID, seatID)
	if err != nil {
		return nil, err
	}
	return order, nil
}