package api

import (
	"chat/api/router"
	"chat/api/rpc"
	"chat/config"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Chat struct {
}

func New() *Chat {
	return &Chat{}
}

func (c *Chat) Run() {
	//初始化rpc客户端
	rpc.InitLogicRpcClient()

	r := router.Register()

	//得到gin的runMode，release或者debug
	//runMode := config.GetGinRunMode()
	apiConfig := config.Conf.Api
	gin.SetMode(apiConfig.ApiBase.GinMode)
	logrus.Info("server start , now run mode is ", apiConfig.ApiBase.GinMode)
	port := apiConfig.ApiBase.ListenPort //监听端口
	flag.Parse()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("start listen : %s\n", err)
		}
	}()
	// if have two quit signal , this signal will priority capture ,also can graceful shutdown
	quit := make(chan os.Signal)
	/*
		SIGHUP：终端控制进程结束(终端连接断开)
		SIGINT：用户发送INTR字符(Ctrl+C)触发
		SIGTERM：结束程序(可以被捕获、阻塞或忽略)
		SIGQUIT：用户发送QUIT字符(Ctrl+/)触发
	*/
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logrus.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//使用Shutdown可以优雅的终止服务，其不会中断活跃连接
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server Shutdown:", err)
	}
	logrus.Infof("Server exiting")
	os.Exit(0)
}
