package schemas

import "super-lender/models"

type GetCustomerContactsResultSchema struct {
	UID         int    `json:"uid"`
	ContactType int    `json:"contactType"`
	Value       string `json:"value"`
	LastUpdate  string `json:"lastUpdate"`
}

type GetCustomerContactResultSchema struct {
	UID         int                `json:"uid" binding:"required,numeric,gt=0"`
	CustomerID  int                `json:"customerId" binding:"required,numeric,gt=0"`
	ContactType models.ContactType `json:"contactType" binding:"required,numeric,oneof=1 2 3"`
	Value       string             `json:"value" binding:"required,min=5,max=50"`
	EncPhone    string             `json:"encPhone" binding:"omitempty"`
	LastUpdate  string             `json:"lastUpdate" binding:"omitempty"`
}
