package config

// GinConfig 定义 Gin 配置文件的结构体
type GinConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// AmapConfig 高德地图配置
type AmapConfig struct {
	APIKey string `mapstructure:"api_key"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	ImageDir string `mapstructure:"image_dir"`
}

// System 定义项目配置文件结构体
type System struct {
	GinConfig    *GinConfig    `mapstructure:"gin"`
	AmapConfig   *AmapConfig   `mapstructure:"amap"`
	UploadConfig *UploadConfig `mapstructure:"upload"`
}
