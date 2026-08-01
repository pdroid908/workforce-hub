package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var app *gin.Engine

func init() {
	// Inisialisasi Gin hanya sekali saat cold start
	gin.SetMode(gin.ReleaseMode)
	app = gin.New()
	app.Use(gin.Recovery())

	// Define Route/Endpoint Gin kamu di sini
	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Halo dari Gin + Vercel!",
			"status":  "success",
		})
	})

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

// Handler adalah fungsi wajib yang diekspor untuk dibaca oleh Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}