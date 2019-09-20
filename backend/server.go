package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func logErrorStackTraceHandler(c *gin.Context) {

}

func logInfoHandler(c *gin.Context) {
	log.Info("This is log with severity of info!")

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func init() {
	// disable debug message
	gin.SetMode(gin.ReleaseMode)

	// setup logger configuration
	// backend log format would be
	// {
	//  "timestamp": ...,
	//  "severity": "info",
	//  "message" : "log message or stacktrace",
	//  "function": "main.logInfoHandler",
	//  "file": "server.go",
	//  }

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "timestamp",
			log.FieldKeyLevel: "severity",
			log.FieldKeyMsg:   "message",
			log.FieldKeyFunc:  "function",
			log.FieldKeyFile:  "file",
		},
	})
	log.SetReportCaller(true)

}

func main() {
	// disable default logger(stdout/stderr)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/info", logInfoHandler)

	router.GET("/errorStackTrace", logErrorStackTraceHandler)

	router.Run(":8080")
}
