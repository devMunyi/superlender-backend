package models

type NextStepStatus uint8

const (
	DeletedNextStep NextStepStatus = iota
	ActiveNextStep  NextStepStatus = 1
)

type ONextStep struct {
	UID     int            `json:"uid" gorm:"primaryKey;autoIncrement"`
	Name    string         `json:"name" gorm:"unique;not null;type:varchar(50)"`
	Details string         `json:"details" gorm:"not null;type:varchar(255)"`
	Status  NextStepStatus `json:"status" gorm:"default:1"`
}
