package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"template/config"
)

var (
	CONFIG      config.System // 系统配置信息
	DB          *gorm.DB
	SugarLogger *zap.SugaredLogger
)

func init() {

}
