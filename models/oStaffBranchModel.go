package models

import "time"

type OStaffBranch struct {
	UID       int       `json:"uid" gorm:"primaryKey;autoIncrement"`
	Agent     int       `json:"agent" gorm:"not null"`
	Branch    int       `json:"branch" gorm:"not null"`
	AddedDate time.Time `json:"added_date" gorm:"autoCreateTime;type:datetime"`
	Status    int       `json:"status" gorm:"default:1"`
}

// TableName specifies the table name for the OStaffBranch model.
func (OStaffBranch) TableName() string {
	return "o_staff_branches"
}
