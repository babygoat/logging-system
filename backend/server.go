package main

import (
	"github.com/gin-gonic/gin"
)

func logErrorStackTraceHandler(c *gin.Context) {

}

func logInfoHandler(c *gin.Context) {

}

func main() {
	// disable default logger(stdout/stderr)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/info", logInfoHandler)

	router.GET("/errorStackTrace", logErrorStackTraceHandler)

	router.Run(":8080")
}
