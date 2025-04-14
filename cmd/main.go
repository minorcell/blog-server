package main

import (
	"demos/internal/controllers"
	"demos/internal/models"
	"demos/internal/services"
	"demos/pkg/middleware"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	router := gin.Default()

	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.Use(gin.Recovery())

	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", userController.LoginUser)
		auth.POST("/register", userController.RegisterUser)
		auth.GET("/user-info", middleware.AuthMiddleware(), userController.GetUserInfo)
		auth.PUT("/user-info", middleware.AuthMiddleware(), userController.UpdateUserInfo)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	serverAddr := fmt.Sprintf(":%s", port)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
