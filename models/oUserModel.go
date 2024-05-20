package models

import "time"

type UserStatus int

const (
	Delete UserStatus = iota
	Active UserStatus = 1
	Blocke UserStatus = 2
)

type OUser struct {
	UID        int        `json:"uid" gorm:"primaryKey;autoIncrement"`
	Name       string     `json:"name" gorm:"type:varchar(50);not null" binding:"required,min=3"`
	Email      string     `json:"email" gorm:"type:varchar(50);not null" binding:"required,email"`
	Phone      string     `json:"phone" gorm:"type:varchar(15);not null;unique" binding:"required,numeric,min=10,max=12"`
	NationalID string     `json:"nationalId" gorm:"type:varchar(15);not null" binding:"omitempty,numeric,min=6"`
	JoinDate   time.Time  `json:"joinDate" gorm:"autoCreateTime;type:datetime"`
	Pass1      string     `json:"password" gorm:"not null" binding:"required,min=6"`
	UserGroup  int        `json:"userGroup" gorm:"not null" binding:"required,numeric,gt=0"`
	Tag        string     `json:"tag" binding:"omitempty,min=2"`
	Pair       int        `json:"pair" gorm:"default:0" binding:"omitempty,numeric,gt=0"`
	Branch     int        `json:"branch" gorm:"not null" binding:"required,numeric,gt=0"`
	Company    int        `json:"company" gorm:"default:1" binding:"omitempty,numeric,gt=0"`
	Status     UserStatus `json:"status" gorm:"default:1" binding:"omitempty,numeric,oneof=0 1 2"`
}
