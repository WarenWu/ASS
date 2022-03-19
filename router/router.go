package router

import (
	"bytes"
	"strings"
	"time"

	"ASS/config"
	//_ "tbsc/docs"
	"ASS/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	//gs "github.com/swaggo/gin-swagger"
	//"github.com/swaggo/gin-swagger/swaggerFiles"
)

var router *gin.Engine

func Router() *gin.Engine {
	return router
}

func InitHttpRouter() {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.New()
	router.Use(AccessLogHandler())
	router.GET("/cn/stockInfos", handler.GetStockInfos)
	router.GET("/cn/stockInfo", handler.GetStockInfo)
	router.POST("/cn/add/stock", handler.AddStock)
	router.POST("/cn/del/stock", handler.DelStock)
	router.GET("/cn/condition", handler.GetStockCondtion)
	router.POST("/cn/set/condition", handler.SetStockCondtion)
	//router.GET("/swagger/*any", gs.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAG"))
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func AccessLogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.RequestURI, "swagger") {
			c.Next()
			return
		}

		if strings.Contains(c.Request.RequestURI, "heartbeat") &&
			(config.LogLevel != "trace" && config.LogLevel != "debug") {
			c.Next()
			return
		}
		blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		startTime := time.Now()

		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logrus.Infof("| %3d | %13v | %15s | %s | %s | %s\n",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
			blw.body.String(),
		)
	}
}