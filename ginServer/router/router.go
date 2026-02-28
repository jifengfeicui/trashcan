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
	router.Use(middle.Cors())

	v1 := router.Group("api")
	{
		v1.GET("/test", api.Test)

		// 用户相关接口（公开）
		v1.POST("/users/register", api.Register)
		v1.POST("/users/login", api.Login)
		v1.POST("/users/wechat-login", api.WechatLogin)

		// 用户相关接口（需要认证）
		authGroup := v1.Group("")
		authGroup.Use(middle.JWTAuth())
		{
			authGroup.GET("/users/me", api.GetCurrentUser)
			authGroup.PUT("/users/me", api.UpdateCurrentUser)
			authGroup.GET("/users/me/trashcans", api.GetUserTrashCans)
		}

		// 垃圾桶相关接口（公开）
		v1.GET("/trashcans/nearby", api.GetNearbyTrashCans)
		v1.GET("/trashcans/:id", api.GetTrashCanDetail)

		// 垃圾桶相关接口（需要认证）
		trashCanAuthGroup := v1.Group("")
		trashCanAuthGroup.Use(middle.JWTAuth())
		{
			trashCanAuthGroup.POST("/trashcans", api.CreateTrashCan)
			trashCanAuthGroup.PUT("/trashcans/:id", api.UpdateTrashCan)
			trashCanAuthGroup.DELETE("/trashcans/:id", api.DeleteTrashCan)
			trashCanAuthGroup.POST("/trashcans/:id/like", api.ToggleLike)
			trashCanAuthGroup.POST("/trashcans/:id/dislike", api.ToggleDislike)
		}
	}

	// 静态文件服务 - 图片访问
	router.Static("/uploads", "./uploads")

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", global.CONFIG.GinConfig.Host, global.CONFIG.GinConfig.Port)

	if global.CONFIG.GinConfig.EnableHTTPS {
		certFile := global.CONFIG.GinConfig.CertFile
		keyFile := global.CONFIG.GinConfig.KeyFile

		if certFile == "" {
			certFile = "cert.pem"
		}
		if keyFile == "" {
			keyFile = "key.pem"
		}

		global.SugarLogger.Infof("Starting HTTPS server on %s", addr)
		global.SugarLogger.Infof("Certificate: %s, Key: %s", certFile, keyFile)
		_ = router.RunTLS(addr, certFile, keyFile)
	} else {
		global.SugarLogger.Infof("Starting HTTP server on %s", addr)
		_ = router.Run(addr)
	}
}
