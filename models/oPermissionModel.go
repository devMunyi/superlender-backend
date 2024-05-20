package models

type OPermission struct {
	UID     int    `json:"uid" gorm:"primaryKey;autoIncrement"`
	GroupID int    `json:"group_id" gorm:"not null"`
	UserID  int    `json:"user_id" gorm:"not null"`
	Tbl     string `json:"tbl" gorm:"type:varchar(50);not null"`
	Rec     int    `json:"rec" gorm:"not null"`
	General int    `json:"general" gorm:"default:0"`
	Create  int    `json:"create" gorm:"default:0"`
	Read    int    `json:"read" gorm:"default:0"`
	Update  int    `json:"update" gorm:"default:0"`
	Delete  int    `json:"delete" gorm:"default:0"`
	Block   int    `json:"block" gorm:"default:0"`
	Unblock int    `json:"unblock" gorm:"default:0"`
}
