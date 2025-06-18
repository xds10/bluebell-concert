package controller

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func CreateConcertHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(models.Concert)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误
		zap.L().Error("create concert with invalid error", zap.Error(err))
		// 判断err是不是validator
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans)) // 翻译
		return
	}
	// 2. 业务逻辑
	if err := logic.CreateConcert(p); err != nil {
		zap.L().Error("logic.create_concert failed", zap.Error(err))
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

// GetConcertListHandler 获取所有的concert
func GetConcertListHandler(c *gin.Context) {
	// 获取数据
	data, err := logic.GetConcertList()
	if err != nil {
		zap.L().Error("logic.GetConcertList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

func GetConcertDetailHandler(c *gin.Context) {
	// 获取id
	idStr := c.Param("id")
	// 类型转换
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 根据id获取数据
	data, err := logic.GetConcertDetail(int64(id))
	if err != nil {
		zap.L().Error("logic.GetConcertDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetConcertSeatsHandler 获取演唱会座位信息
func GetConcertSeatsHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("get concert seats with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取演唱会座位信息
	seatInfo, err := logic.GetConcertSeats(id)
	if err != nil {
		zap.L().Error("logic.GetConcertSeats failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, seatInfo)
}

// GetConcertImageHandler 获取演唱会图片
func GetConcertImageHandler(c *gin.Context) {
	// 1. 获取演唱会ID参数
	idStr := c.Param("id")
	_, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("get concert image with invalid concert id", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 构建图片路径
	imagePath := "./resources/concert-" + idStr + "-picture.jpg" // 默认jpg格式

	// 3. 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		// 尝试png格式
		imagePath = "./resources/concert-" + idStr + "-picture.png"
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			// 图片不存在，返回默认图片或404
			zap.L().Error("concert image not found", zap.String("path", imagePath))
			c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
			return
		}
	}

	// 4. 返回图片文件
	c.File(imagePath)
}
