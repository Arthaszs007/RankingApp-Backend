package router

import (
	"server/controllers"
	"server/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	// config and register logger
	r.Use(gin.LoggerWithConfig(logger.LoggerToFile()))
	r.Use(logger.Recover)

	// config the cors for crossing-origin resource sharing
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://192.168.1.134:5173", "https://ranking-app-frontend.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	// config api
	user := r.Group("/user")
	{
		user.POST("/login", controllers.UserController{}.Login)

		user.POST("/register", controllers.UserController{}.Register)

		user.GET("/verify", controllers.UserController{}.Verify)
	}

	event := r.Group("/event")
	{
		event.GET("list", controllers.EventController{}.GetList)
		event.POST("vote", controllers.EventController{}.VoteToMember)
	}
	rank := r.Group("/rank")
	{
		rank.GET("/list", controllers.RankController{}.GetRank)
	}
	test := r.Group("/test")
	{
		test.GET("1", controllers.RankController{}.GetRank)
	}
	return r
}
