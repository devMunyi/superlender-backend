package models

type OCustomerGuarantorRelationship struct {
	UID    int    `json:"uid" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name" gorm:"type:varchar(50);not null;unique"`
	Status int    `json:"status" gorm:"default:1"`
}

func (OCustomerGuarantorRelationship) TableName() string {
	return "o_customer_guarantor_relationships"
}
