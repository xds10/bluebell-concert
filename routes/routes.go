package routes

import (
	"net/http"

	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 配置CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")

	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.GET("/concert-list", controller.GetConcertListHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) // 认证中间件
	{

		v1.POST("/create_concert", controller.CreateConcertHandler)
		v1.GET("/concert/:id", controller.GetConcertDetailHandler)
		v1.POST("/buy", controller.BuyTicketHandler)
		v1.POST("/pay", controller.PayOrderHandler)
		v1.POST("/order-list", controller.OrderListHandler)
		v1.POST("/order/:id", controller.GetOrderDetailHandler)
		v1.POST("/cancel-order", controller.CancelOrderHandler)
	}

	v1.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return r
}
