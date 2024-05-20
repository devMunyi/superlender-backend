package models

import "time"

type OToken struct {
	UID          uint      `gorm:"primaryKey;autoIncrement"`
	UserID       int       `gorm:"not null"`
	Token        string    `gorm:"type:vachar(245);not null"`
	CreationDate time.Time `gorm:"autoCreateTime;type:datetime"`
	ExpiryDate   time.Time `gorm:"type:datetime"`
	DeviceID     string    `gorm:"type:vachar(245)"`
	BrowserName  string    `gorm:"type:vachar(250)"`
	IPAddress    string    `gorm:"type:vachar(45)"`
	OS           string    `gorm:"type:vachar(55)"`
	Usages       int       `gorm:"default:0"`
	Status       int       `gorm:"default:1;comment:1-valid, 2-expired"`
}
