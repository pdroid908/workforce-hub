package main

// import (
// 	"absen/redis" // Sesuaikan nama module di go.mod kamu
// 	"absen/waqi"  // Import package waqi yang baru dibuat
// 	"fmt"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	_ = godotenv.Load()

// 	_, err := redis.ConnectRedis()
// 	if err != nil {
// 		fmt.Printf("Warning Redis: %v\n", err)
// 	}

// 	r := gin.Default()

// 	// Middleware CORS sederhana agar HTML bisa akses API
// 	r.Use(func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	})
// 	// Add this line to serve index.html directly from Gin
// 	r.StaticFile("/", "./index.html")

// 	// Register handler
// 	r.GET("/api/air-quality", waqi.AirQualityHandler)
// 	r.GET("/api/cities", waqi.GetCitiesDropdownHandler)


// 	r.Run(":8080")
// }
