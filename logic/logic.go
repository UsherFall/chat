package logic

import (
	"chat/config"
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Logic struct {
	ServerId string
}

func New() *Logic {
	return new(Logic)
}

func (logic *Logic) Run() {
	//读取配置文件
	logicConfig := config.Conf.Logic

	//设置当前进程使用的最大cpu数
	runtime.GOMAXPROCS(logicConfig.LogicBase.CpuNum)
	//生成uuid通用唯一识别码
	logic.ServerId = fmt.Sprintf("logic-%s", uuid.New().String())
	//初始化redis客户端
	if err := logic.InitPublishRedisClient(); err != nil {
		logrus.Panicf("logic init publishRedisClient fail,err:%s", err.Error())
	}

	//初始化rpc服务端
	if err := logic.InitRpcServer(); err != nil {
		logrus.Panicf("logic init rpc server fail")
	}
}
