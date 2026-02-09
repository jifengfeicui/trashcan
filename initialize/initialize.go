package initialize

import (
	"template/global"
	Orm "template/initialize/orm"
)

func Initialize() {
	//初始化方法,调试模式勿打开 ChangeWorkingDir
	//ChangeWorkingDir()
	CreateMkdirall()
	InitLogger()
	Viper()
	global.DB = Orm.InitDB()
	Orm.RegisterTables()
}
