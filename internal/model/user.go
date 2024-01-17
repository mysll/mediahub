package model

type User struct {
	ID uint `json:"id" gorm:"primaryKey"` // unique key
}
