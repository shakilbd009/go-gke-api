package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Entry(w http.ResponseWriter, r *http.Response) {

}

func init() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.Run()
}
