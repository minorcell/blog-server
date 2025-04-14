package main

import (
	"fmt"
	"log"
	"os"

	"demos/internal/controllers"
	"demos/internal/models"
	"demos/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(
		&models.User{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	router := gin.Default()

	router.Use(gin.Recovery())

	// 初始化服务层
	userService := services.NewUserService(db)
	// 初始化控制器层
	userController := controllers.NewUserController(userService)
	// 注册路由

	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", userController.LoginUser)
		auth.POST("/register", userController.RegisterUser)
		auth.POST("/logout", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Logout",
			})
		})
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000" // 默认端口
	}

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
