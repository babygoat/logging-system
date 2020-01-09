package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func dbHandler() error {
	_, err := gorm.Open("mysql", "invalid DSN")

	if err != nil {
		return errors.Errorf("Open connection error: %v", err)
	}

	return nil
}

func logErrorStackTraceHandler(c *gin.Context) {
	err := dbHandler()

	if err != nil {
		log.Errorf("%s", FormatStack(err))
	}

	c.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
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
	//  "severity": "info",
	//  "message" : "log message or stacktrace",
	//  "serviceContext" : "{
	//    "service": "poc-backend",
	//    "version": "test"
	//  }"
	//  "context" : "{
	//    "reportLocation": "{
	//       "filePath": "<ERROR FILE PATH>",
	//       "lineNumber": "<EERROR LINE NUMBER>",
	//       "functionName": "ERROR FUNCTION NAME>"
	//    }"
	//  }",
	//  "function": "main.logInfoHandler",
	//  "file": "server.go"
	//  }

	log.SetOutput(os.Stdout)
	log.SetFormatter(NewStackdriverFormatter("poc-backend", "test"))
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
