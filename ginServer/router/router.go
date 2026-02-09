package router

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"template/ginServer/api"
	"template/ginServer/middle"
	"template/global"
)

func InitRouter() {
	router := gin.Default()
	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	router.Use(middle.Cors())
	//router.LoadHTMLFiles("./ui/index.html")
	//router.Static("/static", "./ui/static")
	//router.GET("/", func(c *gin.Context) {
	//	// c.JSON：返回JSON格式的数据
	//	c.HTML(http.StatusOK, "index.html", nil)
	//})

	v1 := router.Group("api")
	{
		v1.GET("/test", api.Test)
		// 垃圾桶相关接口
		v1.GET("/trashcans/nearby", api.GetNearbyTrashCans)
		v1.POST("/trashcans", api.CreateTrashCan)
		v1.GET("/trashcans/:id", api.GetTrashCanDetail)
	}

	// 静态文件服务 - 图片访问
	router.Static("/uploads", "./uploads")
	//ws.RegisterGin(v1)
	//frp.AddFrpRoute(v1)
	_ = router.Run(fmt.Sprintf("%s:%d", global.CONFIG.GinConfig.Host, global.CONFIG.GinConfig.Port))
}
