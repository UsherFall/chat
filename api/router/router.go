package router

import (
	"chat/api/handler"
	"chat/api/rpc"
	"chat/proto"
	"chat/tools"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Register() *gin.Engine {
	r := gin.Default()
	//todo..跨域
	r.Use(CorsMiddleware())
	initUserRouter(r)
	initPushRouter(r)
	//todo..push的路由
	//找不到路由的处理
	r.NoRoute(func(c *gin.Context) {
		tools.FailWithMsg(c, "please check request url !")
	})
	return r
}

func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.POST("/register", handler.Register) //注册
	userGroup.POST("/login", handler.Login)       //登录
	userGroup.Use(CheckSessionId())               //身份认证中间件
	{
		userGroup.POST("/checkAuth", handler.CheckAuth) //用户自主进行身份认证
		userGroup.POST("/logout", handler.Logout)       //登出
	}
}

func initPushRouter(r *gin.Engine) {
	pushGroup := r.Group("/push")
	pushGroup.Use(CheckSessionId())
	{
		pushGroup.POST("/push", handler.Push)               //push单条信息
		pushGroup.POST("/pushRoom", handler.PushRoom)       //push到room中
		pushGroup.POST("/count", handler.Count)             //当前room里的在线人数
		pushGroup.POST("/getRoomInfo", handler.GetRoomInfo) //当前room的用户信息
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		var openCorsFlag = true
		if openCorsFlag {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, nil)
		}
		c.Next()
	}
}

type FormCheckSessionId struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

// 身份认证中间件，如果传来的token错误，则停止后续的请求
func CheckSessionId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formCheckSessionId FormCheckSessionId
		if err := c.ShouldBindBodyWith(&formCheckSessionId, binding.JSON); err != nil {
			//Abort表示终止，也就是说，执行Abort的时候会停止所有的后面的中间件函数的调用
			c.Abort()
			tools.ResponseWithCode(c, tools.CodeSessionError, nil, nil)
			return
		}
		authToken := formCheckSessionId.AuthToken
		req := &proto.CheckAuthRequest{
			AuthToken: authToken,
		}
		code, userId, userName := rpc.RpcLogicObj.CheckAuth(req)
		if code == tools.CodeFail || userId <= 0 || userName == "" {
			c.Abort()
			tools.ResponseWithCode(c, tools.CodeSessionError, nil, nil)
			return
		}
		c.Next()
		return
	}
}
