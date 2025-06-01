package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户存不存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {

		return err
	}
	// 2. 生成UID
	userID := snowflake.GenID()
	//构造user实例
	u := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	// 3. 保存进数据库
	return mysql.InsertUser(u)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递指针,可以拿到userid
	if err = mysql.Login(user); err != nil {
		return user, err
	}
	// 生成JWT
	user.Token, err = jwt.GenToken(user.UserID, user.Username)
	return user, err
}
