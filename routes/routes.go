package routes

import (
	"mta/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/process/:user_id", controllers.Process)
		api.GET("/articles/:user_id", controllers.GetArticles)
		api.GET("/status/:user_id", controllers.GetStatus)
		api.GET("/search/:type/:query", controllers.SearchOnArticles)
	}

	return router
}
