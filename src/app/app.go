package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApp(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}

func StartAppTestLocally() {
	//
}

func init() {
	urlMapping()
	router.Run()
}
