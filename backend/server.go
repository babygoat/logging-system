package main

import (
	"net"
	"net/http"
	"os"
	"strings"

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

func logPanicHandler(c *gin.Context) {
	panic("Ahahhahahhah PANIC!")
}

func logInfoHandler(c *gin.Context) {
	log.Info("This is log with severity of info!")

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Make sure the client closed connection won't trigger
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				log.Errorf("%s", FormatRecover(4))

				if brokenPipe {
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}

			}
		}()
		c.Next()
	}
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
	router.Use(recovery())
	router.Use(gin.LoggerWithFormatter(NewGinLogFormatter()))

	router.GET("/info", logInfoHandler)

	router.GET("/errorStackTrace", logErrorStackTraceHandler)

	router.GET("/panic", logPanicHandler)

	router.Run(":8080")
}
