package ws

import (
	"context"
	"sync"

	"192.168.3.23/resourcelibrary/lsblkjson_parse/model"
	"github.com/lxzan/gws"

	"template/ginServer/api/common"
	"template/utils"
)

// global/websocket.go
var (
	WsManager = NewManager()
	scanOnce  sync.Once
)

type Manager struct {
	// 使用 map 存储连接，Key 可以是用户 ID 或 UUID
	// 使用 RWMutex 保证并发安全
	Clients map[string]*gws.Conn
	Lock    sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		Clients: make(map[string]*gws.Conn),
	}
}

// Add 添加连接
func (m *Manager) Add(id string, conn *gws.Conn) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Clients[id] = conn
}

// Remove 移除连接
func (m *Manager) Remove(id string) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	delete(m.Clients, id)
}

func (m *Manager) Broadcast(msg interface{}) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()

	for _, socket := range m.Clients {
		// 注意：SendWsMsg 内部可能涉及序列化，在大并发下建议在循环外先序列化好
		SendWsMsg(socket, msg.(WsMsg))
	}
}

func startGlobalDiskMonitoring() {
	// 调用你原来的工具函数，传入 Background 保证常驻
	utils.ReadStorageDevicesCtx(context.Background(), func(blockdevices []model.Blockdevice) {
		// [关键] 一旦扫到，直接广播给所有人
		WsManager.Broadcast(WsMsg{
			Code:      common.SUCCESS,
			Data:      blockdevices,
			Msg:       "USB实时更新",
			Operation: "diskinfo",
		})
	})
}
