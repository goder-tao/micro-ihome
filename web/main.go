package main

import (
	"github.com/gin-contrib/sessions"
	sredis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"ihome/web/controller/house"
	"ihome/web/controller/order"
	"ihome/web/controller/user"
	"ihome/web/dao/mysql"
	"ihome/web/dao/redis"
	"ihome/web/utils/middleware"
)

func main() {
	if err := redis.Init(); err != nil {
		log.Error("redis init: ", err)
	}
	mysql.Init()

	r := gin.Default()
	r.Static("/home", "/home/tao/Data/Software/project/go/project/IHome/web/view")

	// session存储在服务端
	store, err := sredis.NewStore(100, "tcp", "localhost:6379", "", []byte("ihome"))
	if err != nil {
		log.Fatalln("session redis store, ", err)
	}
	r.Use(sessions.Sessions("ihome", store))

	api := r.Group("/api")
	{
		// 不需要验证session的路由
		apiv1 := api.Group("/v1.0")
		{
			apiv1.GET("/session", user.GetSession)
			apiv1.GET("/imagecode/:uuid", user.GetImageCd)
			apiv1.GET("/smscode/:mobile", user.GetSmscd)
			apiv1.POST("/users", user.PostRet)
			apiv1.GET("/areas", user.GetArea)
			apiv1.POST("/sessions", user.PostLogin)
			// apiv1.Use(middleware.LoginCheck)
		}

		// 需要验证session的路由
		apiv1Auth := api.Group("/v1.0", middleware.LoginCheck)
		{
			apiv1Auth.DELETE("/session", user.DeleteSession)

			// user路由相关
			apiv1AuthUser := apiv1Auth.Group("/user")
			{
				apiv1AuthUser.GET("", user.GetUserInfo)
				apiv1AuthUser.PUT("/name", user.PutUserName)
				apiv1AuthUser.POST("/avatar", user.PostAvatar)
				apiv1AuthUser.GET("/auth", user.GetUserAuth)
				apiv1AuthUser.POST("/auth", user.PutUserAuth)
				apiv1AuthUser.GET("/houses", user.GetUserHouses)
				apiv1AuthUser.GET("/orders", user.GetUserOrder)
			}

			// houses相关路由
			apiv1AuthHouse := apiv1Auth.Group("/houses")
			{
				apiv1AuthHouse.POST("", house.PostHouses)
				apiv1AuthHouse.GET("", house.GetHouses)
				apiv1AuthHouse.GET("/:id", house.GetHouseInfo)
				apiv1AuthHouse.POST("/:id/images", house.PostHousesImage)
			}
			apiv1Auth.GET("/house/index", house.GetIndex)

			// 订单相关的路由
			apiv1AuthOrder := apiv1Auth.Group("/orders")
			{
				apiv1AuthOrder.GET("", order.PostOrders)
				apiv1AuthOrder.PUT("/:id/status", order.PutOrders)
				apiv1AuthOrder.PUT("/:id/comment", order.PutComment)
			}
		}
	}

	if err := r.Run(":8090"); err != nil {
		log.Fatalln("route run,", err)
	}
}
