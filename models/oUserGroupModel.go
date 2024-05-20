package models

type OUserGroup struct {
	UID         int    `json:"uid" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"type:varchar(30)"`
	Description string `json:"description" gorm:"type:varchar(250)"`
	KPIMeasured int    `json:"kpi_measured" gorm:"default:1"`
	Status      int    `json:"status" gorm:"default:1"`
}
