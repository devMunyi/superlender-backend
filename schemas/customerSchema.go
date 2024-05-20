package schemas

type GetCustomersResultSchema struct {
	UID             int    `json:"uid"`
	PassportPhoto   string `json:"passportPhoto"`
	FullName        string `json:"fullName"`
	Agent           string `json:"agent"`
	EmailAddress    string `json:"emailAddress"`
	PrimaryMobile   string `json:"primaryMobile"`
	Branch          string `json:"branch"`
	PhysicalAddress string `json:"physicalAddress"`
	Status          string `json:"status"`
}

type GetCustomersQueryStringSchema struct {
	PageNo     int    `json:"pageNo" binding:"omitempty,gt=0"`
	PageSize   int    `json:"pageSize" binding:"omitempty,gt=0"`
	OrderBy    string `json:"orderBy" binding:"omitempty,eq=uid"`
	Dir        string `json:"dir" binding:"omitempty,oneof=asc desc"`
	SearchTerm string `json:"searchTerm" binding:"omitempty"`
	CountLimit int    `json:"countLimit" binding:"omitempty,gt=0"`
	Branch     int    `json:"branch" binding:"omitempty,gt=0"`
	Agent      int    `json:"agent" binding:"omitempty,gt=0"`
	Status     int    `json:"status" binding:"omitempty,gt=0"`
}

type GetCustomerResultSchema struct {
	UID              int    `json:"uid"`
	PassportPhoto    string `json:"passportPhoto"`
	FullName         string `json:"fullName"`
	Gender           string `json:"gender"`
	Dob              string `json:"dob"`
	NationalID       string `json:"nationalId"`
	AddedBy          string `json:"addedBy"`
	CurrentCO        string `json:"currentCo"`
	CurrentLO        string `json:"currentLo"`
	LoanLimit        int    `json:"loanLimit"`
	EmailAddress     string `json:"emailAddress"`
	PrimaryMobile    string `json:"primaryMobile"`
	Branch           string `json:"branch"`
	PhysicalAddress  string `json:"physicalAddress"`
	LocationMap      string `json:"locationMap"`
	AddedDate        string `json:"addedDate"`
	Product          string `json:"product"`
	TotalLoans       int    `json:"totalLoans"`
	BadgeIcon        string `json:"badgeIcon"`
	BadgeTitle       string `json:"badgeTitle"`
	BadgeDescription string `json:"badgeDescription"`
	Status           string `json:"status"`
}
