package model

type Test struct {
	ID   int    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name string `json:"name"`
}
