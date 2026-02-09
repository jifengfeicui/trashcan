package model

import "time"

// TrashCan 垃圾桶模型
type TrashCan struct {
	ID          uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	UserID      *uint     `json:"user_id" gorm:"index"` // 可为NULL以兼容现有数据
	Latitude    float64   `json:"latitude" gorm:"type:REAL;not null"`
	Longitude   float64   `json:"longitude" gorm:"type:REAL;not null"`
	Address     string    `json:"address" gorm:"type:TEXT"`
	Description string    `json:"description" gorm:"type:TEXT"`
	ImagePath   string    `json:"image_path" gorm:"type:TEXT"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (TrashCan) TableName() string {
	return "trash_cans"
}
