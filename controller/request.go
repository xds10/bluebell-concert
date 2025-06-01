package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

var ErrorUserNotLogin = errors.New("user not login")

const CtxUserIDKey = "userID"

// GetCurrentUser 获取当前用户id
func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}
func getPageInfo(c *gin.Context) (int64, int64) {
	//获取分页查询
	offsetStr := c.Query("page")
	limitStr := c.Query("size")

	var (
		limit  int64
		offset int64
		err    error
	)
	offset, err = strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 1
	}
	limit, err = strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 5
	}
	return offset, limit
}
