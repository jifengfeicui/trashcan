package model

import "time"

// User 用户模型
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Username     string    `json:"username" gorm:"type:TEXT;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"type:TEXT"`             // 不返回给前端
	OpenID       string    `json:"-" gorm:"type:TEXT;uniqueIndex"` // 微信 openid
	SessionKey   string    `json:"-" gorm:"type:TEXT"`             // 微信 session_key
	Nickname     string    `json:"nickname" gorm:"type:TEXT"`      // 微信昵称
	Avatar       string    `json:"avatar" gorm:"type:TEXT"`        // 微信头像
	WechatID     string    `json:"-" gorm:"type:TEXT"`             // 微信统一ID
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
