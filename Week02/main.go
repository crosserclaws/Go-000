package main

import (
	"github.com/crosserclaws/Go-000/Week02/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/user/:id", api.GetUserInfoByID)
	r.GET("/users", api.GetUsers)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
