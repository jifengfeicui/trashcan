package main

import (
	"trashcan/ginServer"
	"trashcan/global"
	"trashcan/initialize"
)

func main() {
	initialize.Initialize()
	defer global.SugarLogger.Sync()
	global.SugarLogger.Info("start")
	ginServer.Server()
}
