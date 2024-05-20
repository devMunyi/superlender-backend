package models

import "time"

type OEvent struct {
	UID             int        `json:"uid" gorm:"primaryKey;autoIncrement"`
	Tbl             string     `json:"tbl" gorm:"type:varchar(30);not null"`
	Fld             int        `json:"fld" gorm:"not null"`
	EventDetails    string     `json:"event_details" gorm:"type:varchar(250);not null"`
	EventDate       time.Time  `json:"event_date" gorm:"type:datetime;autoCreateTime"`
	EventDateRmTime *time.Time `json:"event_date_rm_time" gorm:"type:date"`
	EventBy         int        `json:"event_by" gorm:"default:0"`
	Status          int        `json:"status" gorm:"default:1"`
	ToSync          int        `json:"to_sync" gorm:"not null;default:0"`
}
