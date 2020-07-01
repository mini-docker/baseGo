package middleware

import (
	"baseGo/src/red_api/app/server"
	"time"

	"baseGo/src/fecho/golog"
)

func TimeLog(next server.HandlerFunc) server.HandlerFunc {
	return func(ctx server.Context) error {
		var t1, t2, t3, t4 int64
		// 程序当前时间
		t1 = time.Now().UnixNano()
		// 进去执行的时间
		err := next(ctx)
		t2 = time.Now().UnixNano()
		// 获取模板的时间
		t3, _ = ctx.Get("tplDownloadTime").(int64)
		// 渲染模板的时间
		t4, _ = ctx.Get("tplExecTime").(int64)
		golog.Info(
			"middleware", "TimeLog",
			"programTotalTime", t2-t1,
			"tplDownloadTime", t3,
			"tplExecTime", t4,
			"RequestURI", ctx.Request().RequestURI)
		return err
	}
}
