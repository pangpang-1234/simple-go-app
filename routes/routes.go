package routes

import (
	"simplegoapp/config"
	"simplegoapp/controllers"
	"simplegoapp/middleware"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {
	db := config.GetDB()
	v1 := r.Group("/api/v1")

	authGroup := v1.Group("auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("/signup", authController.Signup)
		authGroup.POST("/signin", middleware.Authenticate().LoginHandler)
	}
	
	articlesGroup := v1.Group("articles")
	articlesController := controllers.Articles{
		DB: db, // define db config to articles controller
	}
	{
		articlesGroup.GET("", articlesController.FindAll)
		articlesGroup.GET("/:id", articlesController.FindOne)
		articlesGroup.PATCH("/:id", articlesController.Update)
		articlesGroup.DELETE("/:id", articlesController.Delete)
		articlesGroup.POST("", articlesController.Create)
	}

}