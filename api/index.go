package api

import (
	"absen/redis" // Sesuaikan nama module di go.mod kamu
	"absen/waqi"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var app *gin.Engine

func init() {
	_ = godotenv.Load()

	// Konek ke Upstash Redis
	_, err := redis.ConnectRedis()
	if err != nil {
		fmt.Printf("Warning Redis: %v\n", err)
	}

	gin.SetMode(gin.ReleaseMode)
	app = gin.New()
	app.Use(gin.Recovery())

	// Middleware CORS
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Register API Handlers
	app.GET("/api/air-quality", waqi.AirQualityHandler)
	app.GET("/api/cities", waqi.GetCitiesDropdownHandler)
}

// Handler ini diekspor wajib untuk Vercel Serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}