package models

//定义请求参数结构体

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type ParamVoteData struct {
	//从请求中获取当前用户
	PostID    string `json:"post_id" binding:"required"`              //帖子id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` //赞成1 还是 反对-1 取消 0
}

type ParamPostList struct {
	Page  int64  `json:"page" form:"page" `
	Size  int64  `json:"size" form:"size" `
	Order string `json:"order" form:"order" `
}
