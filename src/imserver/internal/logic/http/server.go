package http

import (
	"net"

	"baseGo/src/fecho/echo"
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/middleware"
	"baseGo/src/imserver/internal/logic"
	"baseGo/src/imserver/internal/logic/conf"
)

// Server is http server.
type Server struct {
	engine *echo.Echo
	logic  *logic.Logic
}

// New new a http server.
func New(c *conf.HTTPServer, l *logic.Logic) *Server {
	engine := echo.New()

	s := &Server{
		engine: engine,
		logic:  l,
	}

	s.startServer(c)

	return s
}

func (s *Server) startServer(c *conf.HTTPServer) {
	addr, err := net.ResolveTCPAddr("tcp", ":"+c.Addr)
	if err != nil {
		golog.Error("Server", "startServer", "invalid config address ", err)
		return
	}

	s.engine.Use(middleware.Recover())
	s.initRouter(s.engine)

	err = s.engine.Start(addr.String())
	if err != nil {
		golog.Error("Server", "startServer", "webserver err: %+v", err)
		return
	}
}

func (s *Server) initRouter(e *echo.Echo) {
	group := e.Group("/goim", AuthCheck)
	group.POST("/cachesession", s.CacheSession) //添加session
	group.POST("/joinroom", s.JoinRoom)
	group.POST("/outroom", s.OutRoom)
	group.POST("/delroom", s.DelRoom)
	group.POST("/intiveorkickroom", s.IntiveOrKickRoom)
	group.POST("/push/mids", s.pushMid)                         //	发送消息（个人）
	group.POST("/push/room", s.pushRoom)                        // 发送消息（房间）
	group.POST("/push/group", s.pushGroup)                      // 发送消息（群聊）
	group.POST("/push/room/simulation", s.PushRoomSimulation)   // 发送模拟消息消息（房间）
	group.POST("/push/all", s.pushAll)                          // 发通告
	group.POST("/push/notification", s.sendNotificationMessage) // 发通告
	group.GET("/online/top", s.onlineTop)
	group.GET("/online/room", s.onlineRoom)
	group.GET("/online/total", s.onlineTotal)
	group.GET("/nodes/weighted", s.nodesWeighted)
	group.GET("/nodes/instances", s.nodesInstances)
	group.POST("/online/num", s.RoomsOnlineCount) // 查询房间在线人数
	group.POST("/online/mids", s.RoomsOnlineMids) // 房间mids
}

// Close close the server.
func (s *Server) Close() {

}

func AuthCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return next(ctx)
	}
}
