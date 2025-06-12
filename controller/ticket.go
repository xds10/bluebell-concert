package controller

import (
	"errors"
	"fmt"
	"strconv"

	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func BuyTicketHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.Ticket)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误
		zap.L().Error("buyTicket with invalid error", zap.Error(err))
		// 判断err是不是validator
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans)) // 翻译
		return
	}
	fmt.Println(p)

	// 2. 业务处理
	seatNo, err := logic.BuyTicket(p)
	if err != nil {
		zap.L().Error("logic.BuyTicket failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, seatNo)
}

func PayOrderHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.Order)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误
		zap.L().Error("payTicket with invalid error", zap.Error(err))
		// 判断err是不是validator
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans)) // 翻译
		return
	}
	fmt.Println(p)
	// 2. 业务处理
	err := logic.PayTicket(p)
	if err != nil {
		zap.L().Error("logic.PayTicket failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

func OrderListHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.User)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误
		zap.L().Error("payTicket with invalid error", zap.Error(err))
		// 判断err是不是validator
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans)) // 翻译
		return
	}
	fmt.Println(p)
	// 2. 业务处理
	orders, err := logic.OrderList(p)
	if err != nil {
		zap.L().Error("logic.PayTicket failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, orders)
}

func GetOrderDetailHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	order, err := logic.GetOrderDetail(id)
	if err != nil {
		zap.L().Error("logic.GetOrderDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, order)
}

// 添加取消订单的处理函数
func CancelOrderHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.Order)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误
		zap.L().Error("cancelOrder with invalid error", zap.Error(err))
		// 判断err是不是validator
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans)) // 翻译
		return
	}
	
	// 2. 业务处理
	err := logic.CancelOrder(p.ID)
	if err != nil {
		zap.L().Error("logic.CancelOrder failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
		return
	}
	
	// 3. 返回响应
	ResponseSuccess(c, nil)
}
