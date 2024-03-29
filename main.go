package main

import (
	"InternBCC/database"
	"InternBCC/handler"
	"InternBCC/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	db := database.InitDB()
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to Migrate")
	}
	if err != nil {
		log.Fatalln("failed to load env file")
	}
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	})

	//handler.GDummy()
	//handler.LDummy()
	gedungHandler := handler.NewGedungHandler()
	v1 := r.Group("/v1")
	v1.POST("/register", handler.Register)
	v1.POST("/login", handler.LogIn)
	v1.GET("/validate", middleware.JwtMiddleware(), handler.Validate)
	v1.PUT("/profile/update", middleware.JwtMiddleware(), handler.ChangeNameNumber)
	v1.PUT("/change-pass", middleware.JwtMiddleware(), handler.ChangePass)
	v1.POST("/pass-reset", handler.ForgotPassword)
	v1.DELETE("/delete-account", middleware.JwtMiddleware(), handler.DeleteAccount)

	v1.GET("/gedungs", gedungHandler.FindAllGedung)
	v1.GET("/gedungs/:id", gedungHandler.GetGedungByID)
	v1.GET("/search/gedungs", handler.SearchByName)
	v1.GET("/history", middleware.JwtMiddleware(), handler.GetHistory)
	v1.POST("/booking/:id", middleware.JwtMiddleware(), handler.Booking)
	v1.GET("/booking-details", middleware.JwtMiddleware(), handler.GetBookingData)
	v1.POST("/payment/:id", middleware.JwtMiddleware(), handler.Payment)
	r.Run()
}
