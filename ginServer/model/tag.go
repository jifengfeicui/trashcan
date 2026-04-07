package model

import "time"

// Tag 标签模型
type Tag struct {
	ID        uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name      string     `json:"name" gorm:"type:varchar(100);uniqueIndex;not null"`
	TrashCans []TrashCan `json:"-" gorm:"foreignKey:TagID"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
