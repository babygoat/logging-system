package main

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

func dbHandler() error {
	_, err := gorm.Open("mysql", "invalid DSN")

	if err != nil {
		return xerrors.Errorf("Open connection error: %v", err)
	}

	return nil
}

func logErrorStackTraceHandler(c *gin.Context) {
	err := dbHandler()

	if err != nil {
		log.Errorf("%+v", xerrors.Errorf("db handler error: %v", err))
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
		CallerPrettyfier: func(f *runtime.Frame) (function, file string) {
			function = f.Function
			file = f.File
			return
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
