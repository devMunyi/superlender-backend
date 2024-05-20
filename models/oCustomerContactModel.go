package models

// create a constant for contact types
type ContactType int

const (
	AlternativePhone1 ContactType = 1
	AlternativePhone2 ContactType = 2
	AlternativeEmail1 ContactType = 3
)

type OCustomerContacts struct {
	UID         int     `json:"uid" gorm:"primaryKey;autoIncrement"`
	CustomerID  int     `json:"customer_id" gorm:"not null" binding:"required,numeric,gt=0"`
	ContactType int     `json:"contact_type" gorm:"not null;comment:'From o_contact_types table'" binding:"required,numeric,oneof=1 2 3"`
	Value       string  `json:"value" gorm:"type:varchar(250);not null" binding:"required,max=250"`
	EncPhone    *string `json:"enc_phone" gorm:"type:varchar(70)" binding:"omitempty,max=100"`
	LastUpdate  string  `json:"last_update" gorm:"type:datetime;autoCreateTime" `
	Status      int     `json:"status" gorm:"default:1;comment:'1=Active, 0=Inactive'" binding:"omitempty,numeric,oneof=0 1"`
}

func (OCustomerContacts) TableName() string {
	return "o_customer_contacts"
}
