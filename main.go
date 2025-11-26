package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	server := InitWebServer()

	server.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
