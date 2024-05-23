package schemas

type GetSpCcVintagesResultSchema struct {
	LoanId              int     `json:"loanId"`
	LoanApplicationDate string  `json:"loanApplicationDate"`
	LoanDefaultedDate   string  `json:"loanDefaultedDate"`
	LoanBal             float64 `json:"loanBal"`
	CustomerName        string  `json:"customerName"`
	PhoneNumber         string  `json:"phoneNumber"`
	NationalId          string  `json:"nationalId"`
	Branch              string  `json:"branch"`
	AgentEmail          string  `json:"agentEmail"`
	LoanStatus          string  `json:"loanStatus"`
}
