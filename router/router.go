package router

import (
	"time"

	"gin_gorm_demo/handler"

	"gin_gorm_demo/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"X-Token",
		},
		MaxAge: 12 * time.Hour,
	}))

	api := r.Group("/api/v1")
	{
		api.GET("/users/:id", handler.GetUser)
		api.POST("/users", handler.CreateUser)
		api.PUT("/users/:id", handler.UpdateUser)
		api.DELETE("/users/:id", handler.DeleteUser)
		api.POST("/login", handler.Login)

	}
	api_auth := r.Group("/api/v1/auth")
	api_auth.Use(middleware.AuthMiddleware())
	{
		api_auth.POST("/userlist", handler.GetAuthorizedUsers)
	}

	return r
}
