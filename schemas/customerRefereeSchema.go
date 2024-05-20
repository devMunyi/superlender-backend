package schemas

import "time"

type GetRefereeResultSchema struct {
	UID             int       `json:"uid" gorm:"primary_key;auto_increment"`
	RefereeName     string    `json:"referee_name" gorm:"not null;type:varchar(50)"`
	CustomerId      int       `json:"customer_id" gorm:"not null"`
	MobileNo        string    `json:"mobile_no" gorm:"not null;type:varchar(15);"`
	PhysicalAddress string    `json:"physical_address" gorm:"type:varchar(145)"`
	Relationship    string    `json:"relationship" gorm:"not null"`
	AddedDate       time.Time `json:"added_date" gorm:"not null"`
}

type CreateCustomerRefereeSchema struct {
	RefereeName     string `json:"referee_name" binding:"required,min=3,max=50"`
	CustomerId      int    `json:"customer_id" binding:"required,numeric,gt=0"`
	IdNo            string `json:"id_no" binding:"required,numeric,min=6"`
	MobileNo        string `json:"mobile_no" binding:"required,numeric,min=10,max=12"`
	PhysicalAddress string `json:"physical_address" binding:"omitempty,min=5,max=1000"`
	EmailAddress    string `json:"email_address" binding:"omitempty,email"`
	Relationship    int    `json:"relationship" binding:"required,gt=0"`
	Status          int    `json:"status" binding:"omitempty,numeric,oneof=0 1"`
}

type UpdateCustomerRefereeSchema struct {
	UID             int    `json:"uid" binding:"required,numeric,gt=0"`
	CustomerId      int    `json:"customer_id" binding:"required,numeric,gt=0"`
	AddedDate       string `json:"added_date" binding:"omitempty"`
	RefereeName     string `json:"referee_name" binding:"required,min=3,max=50"`
	IdNo            string `json:"id_no" binding:"required,numeric,min=6"`
	MobileNo        string `json:"mobile_no" binding:"required,numeric,min=10,max=12"`
	PhysicalAddress string `json:"physical_address" binding:"omitempty,min=5,max=1000"`
	EmailAddress    string `json:"email_address" binding:"omitempty,email"`
	Relationship    int    `json:"relationship" binding:"required,gt=0"`
	Status          int    `json:"status" binding:"omitempty,numeric,oneof=0 1"`
}
