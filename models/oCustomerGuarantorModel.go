package models

type OCustomerGuarantor struct {
	UID              int     `json:"uid" gorm:"primaryKey;autoIncrement"`
	GuarantorName    string  `json:"guarantor_name" gorm:"type:varchar(255);not null"`
	CustomerId       int     `json:"customer_id" gorm:"not null"`
	NationalId       string  `json:"national_id" gorm:"type:varchar(15);not null"`
	MobileNo         string  `json:"mobile_no" gorm:"type:varchar(15);not null"`
	PhysicalAddress  string  `json:"physical_address" gorm:"type:text"`
	AmountGuaranteed float64 `json:"amount_guaranteed" gorm:"type:double(100,2);not null;default:0.00"`
	AddedDate        string  `json:"added_date" gorm:"type:datetime;not null"`
	Relationship     int     `json:"relationship" gorm:"not null"`
	Status           int     `json:"status" gorm:"default:1;comment:'1=Active, 0=Inactive'"`
}

func (OCustomerGuarantor) TableName() string {
	return "o_customer_guarantors"
}
