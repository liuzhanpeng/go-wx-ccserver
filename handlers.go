package wxccserver

import "github.com/gin-gonic/gin"

// 返回token
func (s *Server) token(c *gin.Context) {
	appID := c.Query("appid")

	token, err := s.storeManager.Token(appID)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":  0,
		"token": token,
	})
}

// 返回ticket
func (s *Server) ticket(c *gin.Context) {
	appID := c.Query("appid")

	ticket, err := s.storeManager.Ticket(appID)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":   0,
		"ticket": ticket,
	})
}

// 刷新token
func (s *Server) refreshToken(c *gin.Context) {
	appID := c.Query("appid")

	s.storeManager.RemoveToken(appID)

	c.JSON(200, gin.H{
		"code": 0,
	})
}

// 刷新ticket
func (s *Server) refreshTicket(c *gin.Context) {
	appID := c.Query("appid")

	s.storeManager.RemoveTicket(appID)

	c.JSON(200, gin.H{
		"code": 0,
	})
}
