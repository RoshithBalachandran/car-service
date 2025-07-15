package models


type RefreshToken struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint
	Token  string `gorm:"unique"`
}
