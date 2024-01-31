package routes

import (
	"simplegoapp/config"
	"simplegoapp/controllers"

	"github.com/gin-gonic/gin"
)



func Serve(r *gin.Engine) {
	db := config.GetDB()
	articlesGroup := r.Group("/api/v1/articles")
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