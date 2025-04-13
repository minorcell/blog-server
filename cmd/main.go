package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World")

	router := gin.Default()

	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Login",
			})
		})
		auth.POST("/register", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Register",
			})
		})
		auth.POST("/logout", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Logout",
			})
		})
	}

	router.Run(":3000")
}
