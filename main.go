package main

import (
	"template/ginServer"
	"template/global"
	"template/initialize"
)

func main() {
	initialize.Initialize()
	defer global.SugarLogger.Sync()
	global.SugarLogger.Info("start")
	ginServer.Server()
}
