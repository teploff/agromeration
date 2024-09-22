package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHTTPServer(addr string) (*http.Server, error) {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}, nil
}
