package frp

import (
	"192.168.3.23/resourcelibrary/frp_library/client"
	"github.com/gin-gonic/gin"

	"template/global"
)

func AddFrpRoute(rg *gin.RouterGroup) {
	frp := client.New(client.FrpForm{})
	//启动时尝试从文件加载配置（不存在也不致命）
	if err := frp.ReloadFromFile("./frp.json"); err != nil {
		global.SugarLogger.Error("load config failed:", err)
	}
	//注册 frp 控制接口
	client.RegisterGin(rg, frp)
}
