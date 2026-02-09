package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"template/global"
)

var (
	Logger *zap.Logger
	//SugarLogger *zap.SugaredLogger
)

func InitLogger() {
	writeSyncer := getLogWriter()
	consoleWriteSyncer := zapcore.AddSync(os.Stdout) // 添加命令行输出配置
	encoder := getEncoder()
	// 创建多输出的 WriteSyncer
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(writeSyncer, consoleWriteSyncer)
	core := zapcore.NewCore(encoder, multiWriteSyncer, zapcore.DebugLevel)
	//core := zapcore.NewCore(encoder, multiWriteSyncer, zapcore.InfoLevel)
	Logger = zap.New(core, zap.AddCaller())
	global.SugarLogger = Logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 配置显示调用位置（文件名和行号）
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	// 如果想要追加写入可以查看我的博客文件操作那一章
	file, _ := os.Create("./out/stdout.log")
	return zapcore.AddSync(file)
}
