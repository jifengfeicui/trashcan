package utils

import (
	"time"

	"192.168.3.23/resourcelibrary/lsblkjson_parse"
	"192.168.3.23/resourcelibrary/lsblkjson_parse/model"
	"golang.org/x/net/context"

	"template/global"
)

// ReadStorageDevicesCtx 监听信号以中止检测USB设备
func ReadStorageDevicesCtx(ctx context.Context, callback func([]model.Blockdevice)) {
	var intervalTime = time.Duration(2)
	go func() {
		for {
			select {
			case <-ctx.Done():
				global.SugarLogger.Debug("接收到停止信号")
				return
			default:
				blockDevices, err := lsblkjson_parse.ReadForensicDisk()
				if err != nil {
					global.SugarLogger.Error(err)
				}
				if len(blockDevices) > 0 {
					intervalTime = time.Duration(5)
					callback(blockDevices)
				} else {
					intervalTime = time.Duration(2)
				}
			}
			time.Sleep(intervalTime * time.Second)
		}
	}()
}
