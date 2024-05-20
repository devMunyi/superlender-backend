package models

import "time"

type BranchStatus int

const (
	DeletedBranch BranchStatus = iota
	ActiveBranch  BranchStatus = 1
	BlockedBranch BranchStatus = 2
)

type OBranch struct {
	UID                int          `json:"uid" gorm:"primaryKey;autoIncrement"`
	Name               string       `json:"name" gorm:"type:varchar(50)"`
	AddedDate          time.Time    `json:"added_date" gorm:"autoCreateTime;type:datetime"`
	ManagerID          int          `json:"manager_id" gorm:"default:0"`
	AssistantManagerID int          `json:"assistant_manager_id" gorm:"default:0"`
	Address            string       `json:"address" gorm:"type:TEXT;size:1000;not null"`
	RegionID           int          `json:"region_id" gorm:"default:0"`
	Status             BranchStatus `json:"status" gorm:"default:1"`
}

// / specify the table name for the OBranch model.
func (OBranch) TableName() string {
	return "o_branches"
}
