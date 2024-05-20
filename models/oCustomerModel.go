package models

type PhoneNumberProvider int

const (
	SAFARICOM_KE PhoneNumberProvider = 1
	AIRTEL_KE    PhoneNumberProvider = 2
	AIRTEL_UG    PhoneNumberProvider = 3
	MTN_UG       PhoneNumberProvider = 4
)

type Gender string

const (
	MALE   Gender = "M"
	FEMALE Gender = "F"
	OTHER  Gender = "O"
)

type CustomerStatus int // points to customer codes and not uid which starts from 0

const (
	DELETED CustomerStatus = iota
	ACTIVE  CustomerStatus = 1
	BLOCKED CustomerStatus = 2
	LEAD    CustomerStatus = 3
	DRAFT   CustomerStatus = 4

	// Add more Customerstatus options as needed
)

type OCustomer struct {
	UID                 int                 `json:"uid" gorm:"primaryKey;autoIncrement"`
	CustomerCode        string              `json:"customerCode" gorm:"type:varchar(145)"`
	FullName            string              `json:"fullName" gorm:"type:varchar(100);not null" binding:"required,min=3,max=50"`
	PrimaryMobile       string              `json:"primaryMobile" gorm:"type:varchar(15);not null" binding:"required,numeric,min=10,max=12"`
	PhoneNumberProvider PhoneNumberProvider `json:"phoneNumberProvider" gorm:"not null; default:1" binding:"omitempty,numeric,oneof=1 2 3 4"`
	EncPhone            string              `json:"encPhone" gorm:"type:varchar(70);not null"`
	EmailAddress        string              `json:"emailAddress" gorm:"type:varchar(60)" binding:"omitempty,email"`
	PhysicalAddress     string              `json:"physicalAddress" gorm:"type:TEXT;size:1000;not null" binding:"required,min=3,max=1000"`
	Geolocation         string              `json:"geolocation" gorm:"type:TEXT;size:1000" binding:"omitempty,url"`
	Town                int                 `json:"town" gorm:"default:0" binding:"omitempty,numeric"`
	PassportPhoto       string              `json:"passportPhoto" gorm:"type:varchar(255)"`
	NationalID          string              `json:"nationalId" gorm:"type:varchar(10)" binding:"omitempty,numeric,min=6"`
	Gender              Gender              `json:"gender" gorm:"type:varchar(1)" binding:"required,oneof=M F"`
	DOB                 string              `json:"dob" gorm:"type:date"`
	AddedBy             int                 `json:"addedBy" gorm:"not null" binding:"omitempty,numeric"`
	CurrentAgent        int                 `json:"currentAgent" gorm:"default:0" binding:"omitempty,numeric"`
	AddedDate           string              `json:"addedDate" gorm:"autoCreateTime;type:datetime"`
	Branch              int                 `json:"branch" gorm:"not null" binding:"required,numeric,gte=1"`
	PrimaryProduct      int                 `json:"primaryProduct" gorm:"not null" binding:"required,numeric,gte=1"`
	LoanLimit           float64             `json:"loanLimit" gorm:"type:double(100,2);default:0.00" binding:"omitempty,numeric,gte=1"`
	Events              string              `json:"events" gorm:"type:mediumtext" default:""`
	SecData             string              `json:"secData" gorm:"type:longtext"`
	Pin_                string              `json:"pin" gorm:"type:varchar(55)" binding:"omitempty,min=4"`
	DeviceID            string              `json:"deviceId" gorm:"type:varchar(55)" binding:"omitempty,min=4"`
	Flag                int                 `json:"flag" gorm:"default:0" binding:"omitempty,numeric"`
	TotalLoans          int                 `json:"totalLoans" gorm:"default:0"`
	Status              CustomerStatus      `json:"status" gorm:"default:3" binding:"omitempty,numeric,oneof=0 1 2 3 4"`
}
