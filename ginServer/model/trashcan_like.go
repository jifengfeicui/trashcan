package model

import "time"

// TrashCanLike 垃圾桶点赞/点踩模型
type TrashCanLike struct {
	ID         uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	UserID     uint      `json:"user_id" gorm:"uniqueIndex:idx_user_trashcan;not null"`
	TrashCanID uint      `json:"trash_can_id" gorm:"uniqueIndex:idx_user_trashcan;not null"`
	Type       int8      `json:"type" gorm:"type:INTEGER;not null"` // 1=点赞, -1=点踩
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (TrashCanLike) TableName() string {
	return "trash_can_likes"
}

