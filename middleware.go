package wxccserver

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 检查签名
func (s *Server) checkSignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		timestamp, err := strconv.ParseInt(c.Query("timestamp"), 10, 64)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 501,
				"msg":  "无效timestamp",
			})
			c.Abort()
		}

		appID := c.Query("appid")
		signature := c.Query("signature")

		key := s.config.Signature.Key
		if item, ok := s.storeManager.items[appID]; ok {
			if item.signatureKey != "" {
				key = item.signatureKey
			}
		}

		validator := NewHMacSignatureValidator(key, s.config.Signature.Timeout)
		if err := validator.Validate(appID, timestamp, signature); err != nil {

			logrus.Warnf("[%v]签名验证失败:%v", appID, err.Error())

			c.JSON(200, gin.H{
				"code": 501,
				"msg":  "签名验证失败:" + err.Error(),
			})
			c.Abort()
		}
	}
}

// 限制ip访问
func (s *Server) ipLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(s.config.AllowIP) > 0 {
			flag := false
			for _, item := range s.config.AllowIP {
				if item == c.ClientIP() {
					flag = true
				}
			}

			if !flag {
				logrus.Warnf("[%v]IP禁止访问", c.ClientIP())
				c.JSON(200, gin.H{
					"code": 401,
					"msg":  "IP禁止访问",
				})
				c.Abort()
			}
		}
	}
}

// 记录日志
func (s *Server) log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Now().Sub(start)

		logrus.Infof("%v|%v|%v|%v", latency, c.Writer.Status(), c.ClientIP(), c.Request.RequestURI)
	}

}
