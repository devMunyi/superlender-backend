package models

type FlagStatus uint8

const (
	DeletedFlag FlagStatus = iota
	ActiveFlag  FlagStatus = 1
)

type OFlag struct {
	UID         int        `json:"uid" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"unique;not null;type:varchar(50)"`
	Description string     `json:"description" gorm:"not null;type:varchar(255)"`
	ColorCode   string     `json:"color_code" gorm:"not null;type:varchar(10)"`
	Status      FlagStatus `json:"status" gorm:"not null;default:1"`
}
