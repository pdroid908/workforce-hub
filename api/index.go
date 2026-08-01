package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var app *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	app = gin.Default() // Menggunakan Default (sudah include Logger & Recovery)

	// Route utama
	app.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"system":  "WorkforceHub API",
			"status":  "healthy",
			"message": "Enterprise Attendance System API Ready",
		})
	})
}

// Handler wajib diekspor
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}