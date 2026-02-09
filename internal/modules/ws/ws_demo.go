package ws

import (
	"time"

	"github.com/duke-git/lancet/v2/formatter"
	"github.com/lxzan/gws"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"

	"template/ginServer/api/common"
	"template/global"
)

const (
	PingInterval = 5 * time.Second
	PingWait     = 10 * time.Second
)

type Handler struct {
	ID             string
	RootCtx        context.Context
	UsbScanCtx     context.Context // 获取tf卡文件列表时的 cancel
	CancelFuncList []context.CancelFunc
}
type WsMessage struct {
	Operate string      `json:"op"`
	Data    interface{} `json:"data"`
}

func (handler *Handler) InitContext() {
	// 记录 ctx 和 cancelFunc 的映射关系
	handler.RootCtx = context.Background()
	usbScanCtx, usbScanCancel := context.WithCancel(handler.RootCtx)
	handler.UsbScanCtx = usbScanCtx
	handler.CancelFuncList = append(handler.CancelFuncList, usbScanCancel)
}

func (c *Handler) OnOpen(socket *gws.Conn) {
	c.InitContext()
	//_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	//utils.ReadStorageDevicesCtx(c.UsbScanCtx, func(blockdevices []model.Blockdevice) {
	//	SendWsMsg(socket, WsMsg{
	//		Code:      common.SUCCESS,
	//		Data:      blockdevices,
	//		Msg:       "USB信息",
	//		Operation: "diskinfo",
	//	})
	//})
}

func (c *Handler) OnClose(socket *gws.Conn, err error) {
	for _, cancel := range c.CancelFuncList {
		cancel()
	}
	WsManager.Remove(c.ID)
}

func (c *Handler) OnPing(socket *gws.Conn, payload []byte) {
	//_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	_ = socket.WritePong(nil)
}

func (c *Handler) OnPong(socket *gws.Conn, payload []byte) {
}

func (c *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	operationResult := gjson.Get(message.Data.String(), "operation")
	if !operationResult.Exists() {
		global.SugarLogger.Warn("前端未传入指令type")
		SendWsMsg(socket, WsMsg{
			Code:      common.PARAM_EMPTY,
			Data:      nil,
			Msg:       "无效的信息格式",
			Operation: "result",
		})
		return
	}
	switch operationResult.String() {
	case "xxx":
		go wsOperationDemo(socket, message)
	default:
		SendWsMsg(socket, WsMsg{
			Code:      common.PARAM_EMPTY,
			Data:      nil,
			Msg:       "未定义的指令",
			Operation: "result",
		})
	}
}

func initHandler(id string) *gws.Upgrader {
	// 创建 WebSocket 处理器
	handler := &Handler{
		ID: id,
	}
	upgrader := gws.NewUpgrader(handler, &gws.ServerOption{
		ParallelEnabled:   true,                                 // 开启并行消息处理
		Recovery:          gws.Recovery,                         // 开启异常恢复
		PermessageDeflate: gws.PermessageDeflate{Enabled: true}, // 开启压缩
	})
	return upgrader
}

func SendWsMsg(socket *gws.Conn, msg WsMsg) {
	// 检查连接是否有效
	if socket == nil {
		global.SugarLogger.Error("ws错误:", "WebSocket连接无效或已关闭")
		return
	}
	pretty, _ := formatter.Pretty(msg)
	_ = socket.WriteString(pretty)

}

// WsMsg ws的返回信息
// 额外字段用于返回进度信息
type WsMsg struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data,omitempty"`
	Msg       string      `json:"msg"`
	Operation string      `json:"op"`
}

func wsOperationDemo(socket *gws.Conn, message *gws.Message) {
	operation := gjson.Get(message.Data.String(), "operation").String()

	SendWsMsg(socket, WsMsg{
		Code:      common.SUCCESS,
		Data:      "",
		Msg:       "成功",
		Operation: operation,
	})
}
