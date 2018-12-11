package collection

import (
	"net/http"

	ac "grapeGuard/actions"

	log "grapeLoggers/clientv1"

	bl "grapeGuard/blacklist"

	"github.com/labstack/echo"
)

const (
	keepKeys = "43a97140fc0151fc6be19e6f2c267d0147a628b0"
)

func tickOnKeepalive(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{"ret": 0})
}

func tickOnStatus(ctx echo.Context) error {
	keys := ctx.QueryParam("keepAsKey")
	if keys != keepKeys {
		return ctx.JSON(http.StatusOK, echo.Map{
			"ret": 404,
			"msg": "NoFound",
		})
	}

	// 返回数据信息 采集信息
	return ctx.JSON(http.StatusOK, CollectData())
}

func clearCacheData(ctx echo.Context) error {

	var frm ac.ClearCacheNode
	if err := ctx.Bind(&frm); err != nil {
		return ctx.JSON(http.StatusOK, echo.Map{"ret": -1})
	}

	if frm.SignKey != keepKeys {
		return ctx.JSON(http.StatusOK, echo.Map{"ret": 404})
	}

	log.Debug("清理Cache缓存:", frm)

	ac.ClearAction.Push(&frm) //压入队列

	return ctx.JSON(http.StatusOK, echo.Map{"ret": 0})
}

func deletBlackIP(ctx echo.Context) error {
	keys := ctx.QueryParam("delAsKey")
	delIP := ctx.QueryParam("ip")
	if keys != keepKeys {
		return ctx.JSON(http.StatusOK, echo.Map{
			"ret": 404,
			"msg": "NoFound",
		})
	}

	bl.DeleteBlackIP(delIP)
	return ctx.JSON(http.StatusOK, echo.Map{"ret": 0})
}

func SetupHttpCollect(addr string) error {
	e := echo.New()
	e.HideBanner = true
	g := e.Group("/apicoll")
	{
		g.GET("/keepalive", tickOnKeepalive)
		g.GET("/status", tickOnStatus)
		g.GET("/delblackip", deletBlackIP)
		g.POST("/clearCacheData", clearCacheData)
	}

	return e.Start(addr)
}
