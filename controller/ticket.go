package controller

import (
	"errors"
	"fmt"

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
