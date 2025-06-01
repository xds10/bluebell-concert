package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"encoding/hex"
)

const secret = "hello"

// InsertUser 向数据库插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//密码加密
	user.Password, _ = encryptPassword(user.Password)
	// 执行sql语句入库
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

func encryptPassword(opassword string) (string, error) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(opassword))), nil
}
func Login(p *models.User) (err error) {
	opassword, err := encryptPassword(p.Password)
	sqlStr := `select user_id,username,password from user where username=?`
	if err = db.Get(p, sqlStr, p.Username); err != nil {
		return err
	}
	//判断密码是否正确
	if opassword != p.Password {
		return ErrorInvalidPassword
	}
	return nil
}

func GetUserById(id int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select username from user where user_id=?`
	err = db.Get(user, sqlStr, id)
	return
}
