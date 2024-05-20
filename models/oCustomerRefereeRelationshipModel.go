package models

type OCustomerRefereeRelationship struct {
	UID    uint   `json:"uid" gorm:"primary_key;auto_increment"`
	Name   string `json:"name" gorm:"type:varchar(50);unique;not null"`
	Status int    `json:"status" gorm:"default:1"`
}

func (OCustomerRefereeRelationship) TableName() string {
	return "o_customer_referee_relationships"
}
