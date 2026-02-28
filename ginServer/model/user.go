package model

import "time"

// User 用户模型
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Username     string    `json:"username" gorm:"type:TEXT;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"type:TEXT;not null"` // 不返回给前端
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

