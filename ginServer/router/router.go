package router

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"template/ginServer/api"
	"template/ginServer/middle"
	"template/ginServer/static"
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

		// 用户相关接口（公开）
		v1.POST("/users/register", api.Register)
		v1.POST("/users/login", api.Login)

		// 用户相关接口（需要认证）
		authGroup := v1.Group("")
		authGroup.Use(middle.JWTAuth())
		{
			authGroup.GET("/users/me", api.GetCurrentUser)
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

	// 嵌入的前端静态文件服务
	// 创建子文件系统，去掉 dist 前缀
	distFS, err := fs.Sub(static.StaticFiles, "dist")
	if err == nil {
		// 提供静态文件服务（处理所有静态资源：assets、favicon.ico等）
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			// 如果是API请求，返回404
			if len(path) >= 4 && path[:4] == "/api" {
				c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
				return
			}

			// 如果是uploads请求，已经由上面的Static处理了，这里不应该到达
			if len(path) >= 8 && path[:8] == "/uploads" {
				c.Status(http.StatusNotFound)
				return
			}

			// 尝试打开静态文件（去掉前导斜杠）
			filePath := path
			if len(filePath) > 0 && filePath[0] == '/' {
				filePath = filePath[1:]
			}

			// 如果文件路径为空，使用index.html
			if filePath == "" {
				filePath = "index.html"
			}

			file, err := distFS.Open(filePath)
			if err != nil {
				// 文件不存在，返回index.html（SPA路由）
				indexFile, err := distFS.Open("index.html")
				if err != nil {
					c.String(http.StatusInternalServerError, "Failed to load index.html")
					return
				}
				defer indexFile.Close()

				// 读取文件内容
				data, err := io.ReadAll(indexFile)
				if err != nil {
					c.String(http.StatusInternalServerError, "Failed to read index.html")
					return
				}

				// 设置Content-Type
				c.Data(http.StatusOK, "text/html; charset=utf-8", data)
				return
			}
			defer file.Close()

			// 文件存在，读取文件内容
			data, err := io.ReadAll(file)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			// 根据文件扩展名设置Content-Type
			contentType := http.DetectContentType(data)
			if len(filePath) > 0 {
				// 对于已知的文件类型，使用更精确的Content-Type
				if len(filePath) >= 5 && filePath[len(filePath)-5:] == ".html" {
					contentType = "text/html; charset=utf-8"
				} else if len(filePath) >= 4 && filePath[len(filePath)-4:] == ".css" {
					contentType = "text/css; charset=utf-8"
				} else if len(filePath) >= 3 && filePath[len(filePath)-3:] == ".js" {
					contentType = "application/javascript; charset=utf-8"
				} else if len(filePath) >= 4 && filePath[len(filePath)-4:] == ".ico" {
					contentType = "image/x-icon"
				}
			}

			c.Data(http.StatusOK, contentType, data)
		})
	}

	//ws.RegisterGin(v1)
	//frp.AddFrpRoute(v1)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", global.CONFIG.GinConfig.Host, global.CONFIG.GinConfig.Port)

	if global.CONFIG.GinConfig.EnableHTTPS {
		// HTTPS 模式
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
		// HTTP 模式
		global.SugarLogger.Infof("Starting HTTP server on %s", addr)
		_ = router.Run(addr)
	}
}
