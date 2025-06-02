package routes

import (
	"net/http"

	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	v1 := r.Group("/api/v1")

	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.GET("/concert-list", controller.GetConcertListHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) // 认证中间件
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts", controller.GetPostListHandler)
		v1.GET("/posts2", controller.GetPostListHandler2)
		// 投票
		v1.POST("/vote", controller.PostVoteController)

		v1.POST("/create_concert", controller.CreateConcertHandler)
		v1.GET("/concert/:id", controller.GetConcertDetailHandler)
		v1.POST("/buy", controller.BuyTicketHandler)
		v1.POST("/pay", controller.PayOrderHandler)
	}

	v1.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return r
}
