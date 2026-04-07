package model

import "time"

// TrashCan 垃圾桶模型
type TrashCan struct {
	ID          uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	UserID      *uint     `json:"user_id" gorm:"index"`        // 可为NULL以兼容现有数据
	User        *User     `json:"-" gorm:"foreignKey:UserID"`  // 关联用户
	TagID       *uint     `json:"tag_id" gorm:"index"`         // 关联标签
	Tag         *Tag      `json:"tag" gorm:"foreignKey:TagID"` // 预加载标签
	Latitude    float64   `json:"latitude" gorm:"type:REAL;not null"`
	Longitude   float64   `json:"longitude" gorm:"type:REAL;not null"`
	Address     string    `json:"address" gorm:"type:TEXT"`
	Description string    `json:"description" gorm:"type:TEXT"`
	ImagePath   string    `json:"image_path" gorm:"type:TEXT"`
	ImagePath2  string    `json:"image_path_2" gorm:"type:TEXT"`
	ImagePath3  string    `json:"image_path_3" gorm:"type:TEXT"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (TrashCan) TableName() string {
	return "trash_cans"
}
