package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// SaveImage 保存上传的图片文件
// file: 上传的文件
// uploadDir: 上传目录（相对于项目根目录）
// 返回：文件路径（相对于项目根目录）和错误
func SaveImage(file *multipart.FileHeader, uploadDir string) (string, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 检查文件类型
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("文件类型错误，只支持图片格式")
	}

	// 检查文件大小（限制为10MB）
	if file.Size > 10*1024*1024 {
		return "", fmt.Errorf("文件大小超过限制（10MB）")
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 构建完整文件路径
	filePath := filepath.Join(uploadDir, filename)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	// 返回相对路径（用于数据库存储）
	return filePath, nil
}

// GetImageURL 获取图片的访问URL
// imagePath: 图片路径（数据库存储的路径）
// 返回：图片的URL路径
func GetImageURL(imagePath string) string {
	if imagePath == "" {
		return ""
	}
	// 将路径中的反斜杠转换为正斜杠，用于URL
	urlPath := strings.ReplaceAll(imagePath, "\\", "/")
	// 如果路径不是以/开头，添加/
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	return urlPath
}

// EnsureUploadDir 确保上传目录存在
func EnsureUploadDir(uploadDir string) error {
	return os.MkdirAll(uploadDir, 0755)
}
