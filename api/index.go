package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var app *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	app = gin.New()
	app.Use(gin.Recovery())

	app.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"system":  "WorkforceHub API",
			"status":  "healthy",
			"message": "Enterprise Attendance System API Ready",
		})
	})
}

// ⚠️ NAMA FUNGSI HARUS 'Handler' (H Kapital) dan MENGGUNAKAN 'http.ResponseWriter', '*http.Request'
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}