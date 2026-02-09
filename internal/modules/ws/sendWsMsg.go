package ws

//func SendWsMsg(socket *gws.Conn, msg WsMsg) {
//	// 检查连接是否有效
//	if socket == nil {
//		global.SugarLogger.Error("ws错误:", "WebSocket连接无效或已关闭")
//		return
//	}
//	//pretty, _ := formatter.Pretty(msg)
//	//_ = socket.WriteString(pretty)
//	data, err := json.Marshal(msg)
//	if err != nil {
//		// 建议增加错误处理，避免发送空数据或破碎数据
//		global.SugarLogger.Error("JSON marshal failed: %v", err)
//		return
//	}
//
//	// 直接调用 WriteMessage 传入 []byte，效率最高
//	err = socket.WriteString(string(data))
//	global.SugarLogger.Error(err)
//}
