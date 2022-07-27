package protocol

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/yuanyp8/cmdb/apps"
	"github.com/yuanyp8/cmdb/conf"
	"net/http"
	"time"
)

type HttpService struct {
	server *http.Server
	l      logger.Logger
	r      gin.IRouter
}

func NewHttpService() *HttpService {
	r := gin.Default()

	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.HTTPAddr(),
		Handler:           r,
	}

	return &HttpService{
		server: server,
		l:      zap.L().Named("HTTP Service"),
		r:      r,
	}
}

func (s *HttpService) Start() error {
	// 加载IOC层的GinAPPS
	apps.InitGin(s.r)

	// 打印已加载App的日志信息
	appSlice := apps.LoadedGinApps()
	s.l.Infof("loaded gin apps :%v", appSlice)

	// 这个操作是阻塞的，会监听中断信号
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("service stoped success")
			return nil
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}
	return nil
}

func (s *HttpService) Stop() {
	s.l.Info("start graceful shutdown")
	// 给一个关闭的窗口时间，并设置超时时间，如超时未处理完成则强制关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Warnf("shut down http service error, %s", err)
	}
}
