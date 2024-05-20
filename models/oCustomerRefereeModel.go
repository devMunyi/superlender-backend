package models

import "time"

type OCustomerReferee struct {
	UID             int       `json:"uid" binding:"omitempty,numeric,gt=0" gorm:"primary_key;auto_increment"`
	CustomerId      int       `json:"customer_id" binding:"required,numeric,gt=0" gorm:"not null"`
	AddedDate       time.Time `json:"added_date" gorm:"not null" binding:"omitempty"`
	RefereeName     string    `json:"referee_name" gorm:"not null;type:varchar(50)" binding:"required,min=3,max=50"`
	IdNo            string    `json:"id_no" gorm:"not null;type:varchar(15)" binding:"omitempty,numeric,min=6"`
	MobileNo        string    `json:"mobile_no" gorm:"not null;type:varchar(15)" binding:"required,numeric,min=10,max=12"`
	PhysicalAddress string    `json:"physical_address" gorm:"type:varchar(145)" binding:"omitempty,min=5,max=1000"`
	EmailAddress    string    `json:"email_address" gorm:"type:varchar(50)" binding:"omitempty,email"`
	Relationship    int       `json:"relationship" gorm:"not null" binding:"required,gt=0"`
	Status          int       `json:"status" gorm:"default:1" binding:"omitempty,numeric,oneof=0 1"`
}

func (OCustomerReferee) TableName() string {
	return "o_customer_referees"
}
