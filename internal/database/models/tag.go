package models

type Tag struct {
	ID int `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
}
