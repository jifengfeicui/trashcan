package utils

import (
	"github.com/google/uuid"
)

// GenerateUUID 生成一个随机的 V4 UUID 字符串
func GenerateUUID() string {
	return uuid.New().String()
}
