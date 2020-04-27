package wxccserver

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Server 中控服务器
type Server struct {
	config       *Config       // 服务器配置
	storeManager *StoreManager // token/ticket存储管理器
	srv          *http.Server  // http服务器
}

// NewServer 创建中控服务器
func NewServer(config *Config) *Server {
	storeManger := NewStoreManager(NewTomlAccountProvider(config.AccountsFile), NewWxFetcher(), config.AheadTimeout)
	storeManger.LoadAccounts()

	return &Server{
		config:       config,
		storeManager: storeManger,
	}
}

// Start 启动服务
func (s *Server) Start() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(s.checkSignature(), s.ipLimit(), s.log(), gin.Recovery())

	router.GET("/token", s.token)
	router.GET("/ticket", s.ticket)
	router.GET("/refresh-token", s.refreshToken)
	router.GET("/refresh-ticket", s.refreshTicket)

	s.srv = &http.Server{
		Addr:    ":" + strconv.Itoa(s.config.Port),
		Handler: router,
	}

	if s.config.EnableTLS {
		if err := s.srv.ListenAndServeTLS(s.config.TLSCertFile, s.config.TLSKeyFile); err != nil && err != http.ErrServerClosed {
			logrus.Panicf("服务监听失败:%v", err)
		}
	} else {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Panicf("服务监听失败:%v", err)
		}
	}
}

// Stop 停止服务
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		logrus.Panicf("服务停止失败: %v", err)
	}
}

// Reload 重新加载公众号账号信息
func (s *Server) Reload() error {
	return s.storeManager.LoadAccounts()
}
