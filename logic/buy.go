package logic

import (
	"fmt"
	"time"
	"sync"

	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"

	"go.uber.org/zap"
)

// 用于跟踪订单处理状态的全局映射
var (
	// 存储订单处理状态的映射: key=concertID:seatID, value=orderID
	pendingOrders     = make(map[string]int64)
	pendingOrdersLock sync.RWMutex
)

// 获取订单ID
func GetOrderIDBySeat(concertID, seatID int64) (int64, bool) {
	key := fmt.Sprintf("%d:%d", concertID, seatID)
	pendingOrdersLock.RLock()
	defer pendingOrdersLock.RUnlock()
	orderID, exists := pendingOrders[key]
	return orderID, exists
}

// 设置订单ID
func SetOrderIDBySeat(concertID, seatID, orderID int64) {
	key := fmt.Sprintf("%d:%d", concertID, seatID)
	pendingOrdersLock.Lock()
	defer pendingOrdersLock.Unlock()
	pendingOrders[key] = orderID
}

func BuyTicket(p *models.Ticket) (*models.Ticket, error) {
	// 1. 从Redis获取座位
	seatID, err := redis.GetSeatByConcertID(p)
	if seatID == 0 {
		zap.L().Error("redis.GetSeatByConcertID failed", zap.Error(err))
		return nil, fmt.Errorf("获取座位失败: %w", err)
	}
	zap.L().Info("获取到座位", zap.Int64("座位ID", seatID))
	
	if p.SeatIdx == nil {
		p.SeatIdx = &models.Seat{}
	}
	p.SeatIdx.SeatID = seatID
	
	// 准备返回值 - 只包含座位信息
	result := &models.Ticket{
		ConcertID: p.ConcertID,
		UserID:    p.UserID,
		SeatIdx: &models.Seat{
			SeatID: seatID,
		},
	}
	
	// 2. 异步处理数据库操作
	go func() {
		// 开始数据库事务
		tx, err := mysql.DB().Beginx()
		if err != nil {
			zap.L().Error("begin transaction failed", zap.Error(err))
			return
		}
		
		// 声明事务结果变量
		var txErr error
		
		// 使用defer处理事务提交或回滚
		defer func() {
			if p := recover(); p != nil {
				// 发生恐慌时回滚事务
				_ = tx.Rollback()
				zap.L().Error("购票协程异常", zap.Any("panic", p))
			} else if txErr != nil {
				// 发生错误时回滚事务并记录日志
				_ = tx.Rollback()
				zap.L().Error("购票事务失败", zap.Error(txErr))
			}
		}()
		
		// 3. 在事务中执行购票操作
		if txErr = mysql.BuyTicketTx(tx, p); txErr != nil {
			zap.L().Error("mysql.BuyTicket failed", zap.Error(txErr))
			return
		}
		
		// 4. 在事务中查询完整的seat表
		p.SeatIdx, txErr = mysql.GetSeatByIDTx(tx, seatID)
		if txErr != nil {
			zap.L().Error("mysql.GetSeatByID failed", zap.Error(txErr))
			return
		}
		
		// 5. 在事务中创建订单
		order := &models.Order{
			UserID:     p.UserID,
			ConcertID:  p.ConcertID,
			SeatID:     seatID,
			Price:      float64(p.SeatIdx.Price),
			Status:     1,
			CreateTime: time.Now(),
		}
		
		txErr = mysql.AddOrderTx(tx, order)
		if txErr != nil {
			zap.L().Error("mysql.AddOrders failed", zap.Error(txErr))
			return
		}
		
		// 提交事务
		if txErr = tx.Commit(); txErr != nil {
			zap.L().Error("commit transaction failed", zap.Error(txErr))
			return
		}
		
		// 成功完成事务后，将订单ID存入映射表
		SetOrderIDBySeat(p.ConcertID, seatID, order.ID)
		zap.L().Info("异步创建订单成功", 
			zap.Int64("concertID", p.ConcertID), 
			zap.Int64("seatID", seatID),
			zap.Int64("orderID", order.ID))
	}()
	
	return result, nil
}

// 根据演唱会ID和座位ID查询订单
func GetOrderByTicketInfo(concertID, seatID int64) (*models.Order, error) {
	// 首先检查内存中是否有该订单的处理结果
	if orderID, exists := GetOrderIDBySeat(concertID, seatID); exists {
		return mysql.GetOrderById(orderID)
	}
	
	// 如果内存中没有，查询数据库
	order, err := mysql.GetOrderBySeat(concertID, seatID)
	if err != nil {
		return nil, err
	}
	
	// 将查询到的订单ID缓存到内存
	if order != nil {
		SetOrderIDBySeat(concertID, seatID, order.ID)
	}
	
	return order, nil
}

func PayTicket(p *models.Order) error {
	order, err := mysql.GetOrderById(p.ID)
	if (time.Now().Sub(order.CreateTime)).Minutes() > 5 {
		zap.L().Error("订单超时")
		return fmt.Errorf("订单超时")
	}
	err = mysql.PayOrder(p.ID)
	if err != nil {
		zap.L().Error("mysql.PayTicket failed", zap.Error(err))
		return err
	}
	return nil
}

func OrderList(user *models.User) (data []*models.Order, err error) {
	orders, err := mysql.GetOrderList(user.UserID)
	if err != nil {
		zap.L().Error("mysql.GetOrderList failed", zap.Error(err))
		return nil, err
	}
	return orders, nil
}

func GetOrderDetail(p int) (*models.Order, error) {
	order, err := mysql.GetOrderById(int64(p))
	if order == nil {
		return nil, fmt.Errorf("订单不存在")
	}
	order.Ticket = &models.Ticket{
		TicketID:  order.ID,
		ConcertID: order.ConcertID,
		UserID:    order.UserID,
	}
	order.Ticket, err = mysql.GetTicketByOrderID(order.SeatID)
	if err != nil {
		zap.L().Error("mysql.GetTicketByOrderID failed", zap.Error(err))
		return nil, err
	}
	return order, nil
}
