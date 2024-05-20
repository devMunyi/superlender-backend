package models

type OGuarantor struct {
	UID              int     `json:"uid" gorm:"primaryKey;autoIncrement"`
	CustomerID       int     `json:"customer_id" gorm:"not null"`
	NationalID       string  `json:"national_id" gorm:"type:varchar(15);not null"`
	MobileNo         string  `json:"mobile_no" gorm:"type:varchar(15);not null"`
	AmountGuaranteed float64 `json:"amount_guaranteed" gorm:"type:double(100,2);not null"`
	AddedDate        string  `json:"added_date" gorm:"type:datetime;not null"`
	Status           int     `json:"status" gorm:"default:1;comment:'1=Active, 0=Inactive'"`
}

func (OGuarantor) TableName() string {
	return "o_guarantors"
}
