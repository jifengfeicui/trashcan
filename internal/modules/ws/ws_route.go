package ws

import (
	"github.com/gin-gonic/gin"

	"template/utils"
)

func RegisterGin(rg *gin.RouterGroup) {
	rg.GET("/ws", WsDemo)
}
func WsDemo(c *gin.Context) {
	scanOnce.Do(func() {
		go startGlobalDiskMonitoring()
	})

	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = utils.GenerateUUID()
	}
	upgrader := initHandler(clientID)
	socket, err := upgrader.Upgrade(c.Writer, c.Request)
	if err != nil {
		return
	}
	WsManager.Add(clientID, socket)
	go func() {
		socket.ReadLoop() // 此处阻塞会使请求上下文不能顺利被GC
	}()
}
