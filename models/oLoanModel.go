package models

import (
	"time"
)

type LoanStatus uint8

const (
	Created       LoanStatus = 1
	Pending       LoanStatus = 2
	Disbursed     LoanStatus = 3
	PartiallyPaid LoanStatus = 4
	Cleared       LoanStatus = 5
	Rejected      LoanStatus = 6
	Overdue       LoanStatus = 7
	MissedPayment LoanStatus = 8
	WriteOff      LoanStatus = 9
	WrittenOff    LoanStatus = 10
	Reversed      LoanStatus = 11
)

type LoanApplicationMode string

const (
	Manual LoanApplicationMode = "MANUAL"
	USSD   LoanApplicationMode = "USSD"
	SMS    LoanApplicationMode = "SMS"
	APP    LoanApplicationMode = "APP"
)

type OLoan struct {
	UID                     int                 `json:"uid" gorm:"primaryKey;autoIncrement"`
	LoanCode                string              `json:"loan_code" gorm:"type:varchar(50)"`
	CustomerID              int                 `json:"customer_id" gorm:"not null"`
	GroupID                 int                 `json:"group_id" gorm:"default:0"`
	AccountNumber           string              `json:"account_number" gorm:"type:varchar(30)"`
	EncPhone                string              `json:"enc_phone" gorm:"type:varchar(70)"`
	ProductID               int                 `json:"product_id" gorm:"not null"`
	LoanType                int                 `json:"loan_type" gorm:"default:0"`
	LoanAmount              float64             `json:"loan_amount" gorm:"type:double(50,2);not null"`
	DisbursedAmount         float64             `json:"disbursed_amount" gorm:"type:double(50,2);not null"`
	TotalRepayableAmount    float64             `json:"total_repayable_amount" gorm:"type:double(50,2);not null"`
	TotalRepaid             float64             `json:"total_repaid" gorm:"type:double(50,2);default:0.00"`
	LoanBalance             float64             `json:"loan_balance" gorm:"type:double(50,2);not null"`
	Period                  int                 `json:"period" gorm:"not null"`
	PeriodUnits             string              `json:"period_units" gorm:"type:varchar(30)"`
	PaymentFrequency        string              `json:"payment_frequency" gorm:"type:varchar(30)"`
	PaymentBreakdown        string              `json:"payment_breakdown" gorm:"type:varchar(75)"`
	TotalAddons             float64             `json:"total_addons" gorm:"type:double(50,2);not null"`
	TotalDeductions         float64             `json:"total_deductions" gorm:"type:double(50,2);default:0.00"`
	TotalInstalments        int                 `json:"total_instalments" gorm:"not null"`
	TotalInstalmentsPaid    int                 `json:"total_instalments_paid" gorm:"default:0"`
	CurrentInstalment       int                 `json:"current_instalment" gorm:"default:1"`
	CurrentInstalmentAmount float64             `json:"current_instalment_amount" gorm:"type:double(50,2);not null"`
	IncomeEarned            float64             `json:"income_earned" gorm:"default:0.00"`
	GivenDate               time.Time           `json:"given_date" validate:"required" gorm:"type:date;not null"`
	NextDueDate             time.Time           `json:"next_due_date" gorm:"type:date;not null"`
	FinalDueDate            time.Time           `json:"final_due_date" gorm:"type:date;not null"`
	LastPayDate             *time.Time          `json:"last_pay_date" gorm:"type:date"`
	AddedBy                 int                 `json:"added_by" gorm:"default:0"`
	CurrentAgent            int                 `json:"current_agent" gorm:"default:0"`
	CurrentLO               int                 `json:"current_lo" gorm:"default:0"`
	CurrentCO               int                 `json:"current_co" gorm:"default:0"`
	Allocation              string              `json:"allocation" gorm:"default:BRANCH"`
	CurrentBranch           int                 `json:"current_branch" gorm:"not null"`
	AddedDate               time.Time           `json:"added_date" gorm:"autoCreateTime;type:datetime"`
	LoanStage               int                 `json:"loan_stage" gorm:"not null"`
	LoanFlag                int                 `json:"loan_flag" gorm:"default:0"`
	PendingEvent            int                 `json:"pending_event" gorm:"default:0"`
	TransactionCode         string              `json:"transaction_code"`
	TransactionDate         time.Time           `json:"transaction_date" gorm:"type:datetime"`
	ApplicationMode         LoanApplicationMode `json:"application_mode" gorm:"type:varchar(40)"`
	DisburseState           string              `json:"disburse_state" gorm:"type:varchar(45);default:NONE"`
	Disbursed               int                 `json:"disbursed" gorm:"default:0"`
	Paid                    int                 `json:"paid" gorm:"default:0"`
	OtherInfo               string              `json:"other_info" gorm:"type:longtext"`
	Status                  LoanStatus          `json:"status" gorm:"default:1"`
	ToSync                  int                 `json:"to_sync" gorm:"default:0"`
}
