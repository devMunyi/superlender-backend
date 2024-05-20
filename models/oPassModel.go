package models

type OPass struct {
	UID            int    `json:"uid" gorm:"primaryKey;autoIncrement"`
	User           int    `json:"user" gorm:"unique;not null"`
	Pass           string `json:"pass" gorm:"type:varchar(200);not null"`
	PassResetToken string `json:"pass_reset_token" gorm:"type:varchar(255)"`
	ResetStatus    int    `json:"reset_status" gorm:"default:0"`
}

// TableName specifies the table name for the Pass struct
func (OPass) TableName() string {
	return "o_passes"
}
