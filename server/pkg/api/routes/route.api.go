package routes

import (
	"net/http"

	"github.com/emarifer/calorie-tracker-app/pkg/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutesApi(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Check API
		api.GET("/healthchecker", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Hello from Gin!!",
			})
		})

		// API Routes
		api.POST("/entry/create", handlers.AddEntry)
		api.GET("/entries", handlers.GetEntries)
		api.GET("/entry/:id", handlers.GetEntryById)
		api.GET("/ingredient/:ingredient", handlers.GetEntriesByIngredient)

		api.PUT("/entry/update/:id", handlers.UpdateEntry)
		api.PUT("/ingredient/update/:id", handlers.UpadateIngredient)
		api.DELETE("/entry/:id", handlers.DeleteEntry)
	}
}

/* Grouping en Gin. VER:
https://gin-gonic.com/es/docs/examples/grouping-routes/
*/
