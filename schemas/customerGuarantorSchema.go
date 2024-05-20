package schemas

type CreateCustomerGuarantorSchema struct {
	GuarantorName    string  `json:"guarantor_name" binding:"required,min=3,max=50"`
	CustomerId       int     `json:"customer_id" binding:"required,numeric,gt=0"`
	MobileNo         string  `json:"mobile_no" binding:"required,numeric,min=10,max=12"`
	NationalId       string  `json:"national_id" binding:"required,numeric,min=6"`
	PhysicalAddress  string  `json:"physical_address" binding:"required,min=5,max=1000"`
	AmountGuaranteed float64 `json:"amount_guaranteed" binding:"omitempty,numeric,gte=0"`
	AddedDate        string  `json:"added_date" binding:"required,min=19"`
	Relationship     int     `json:"relationship" binding:"required,gt=0"`
	Status           int     `json:"status" binding:"omitempty,numeric,oneof=0 1"`
}

// customerGuarantor := models.OCustomerGuarantor{
// 	UID:              customerGuarantorInput.UID,
// 	GuarantorName:    guarantorName,
// 	CustomerID:       customerGuarantorInput.CustomerId,
// 	NationalID:       nationalId,
// 	MobileNo:         mobileNo,
// 	PhysicalAddress:  physicalAddress,
// 	AmountGuaranteed: customerGuarantorInput.AmountGuaranteed,
// 	Relationship:     customerGuarantorInput.Relationship,
// 	Status:           customerGuarantorInput.Status,
// }

type UpdateCustomerGuarantorSchema struct {
	UID              int     `json:"uid" binding:"required,numeric,gt=0"`
	GuarantorName    string  `json:"guarantor_name" binding:"required,min=3,max=50"`
	CustomerId       int     `json:"customer_id" binding:"required,numeric,gt=0"`
	NationalId       string  `json:"national_id" binding:"required,numeric,min=6"`
	MobileNo         string  `json:"mobile_no" binding:"required,numeric,min=10,max=12"`
	PhysicalAddress  string  `json:"physical_address" binding:"required,min=5,max=1000"`
	AmountGuaranteed float64 `json:"amount_guaranteed" binding:"omitempty,numeric,gte=0"`
	AddedDate        string  `json:"added_date" binding:"required,min=19"`
	Relationship     int     `json:"relationship" binding:"required,gt=0"`
	Status           int     `json:"status" binding:"omitempty,numeric,oneof=0 1"`
}

type GetCustomerGuarantorsResultSchema struct {
	UID              int     `json:"uid"`
	GuarantorName    string  `json:"guarantor_name"`
	CustomerID       int     `json:"customer_id"`
	MobileNo         string  `json:"mobile_no"`
	NationalID       string  `json:"national_id"`
	PhysicalAddress  string  `json:"physical_address"`
	AmountGuaranteed float64 `json:"amount_guaranteed"`
	AddedDate        string  `json:"added_date"`
	Relationship     string  `json:"relationship"`
	Status           int     `json:"status"`
}

type GetCustomerGuarantorResultSchema struct {
	GuarantorName    string  `json:"guarantor_name" binding:"required,min=3,max=50"`
	CustomerId       int     `json:"customer_id" binding:"required,numeric,gt=0"`
	NationalId       string  `json:"national_id" binding:"required,numeric,min=6"`
	MobileNo         string  `json:"mobile_no" binding:"required,numeric,min=10,max=12"`
	PhysicalAddress  string  `json:"physical_address" binding:"required,min=5,max=1000"`
	AmountGuaranteed float64 `json:"amount_guaranteed" binding:"omitempty,numeric,gte=0"`
	AddedDate        string  `json:"added_date" binding:"required,min=19"`
	Relationship     int     `json:"relationship" binding:"required,gt=0"`
}
