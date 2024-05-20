package models

type PlatformSetting struct {
	UID       uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(50);not null"`
	CompanyID int    `gorm:"default:0"`
	Logo      string `gorm:"type:mediumtext;not null"`
	Icon      string `gorm:"type:mediumtext;not null"`
	Link      string `gorm:"type:mediumtext;not null"`
}
