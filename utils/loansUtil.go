package utils

import "gorm.io/gorm"

func FindCCVintages(db *gorm.DB, queryType string) *gorm.DB {
	query := db.Table("o_loans l")

	if queryType == "count" {
		query = query.Select("l.uid")
	} else {
		query = query.Select("l.uid AS loanId, l.give_date AS loanApplicationDate, l.final_due_date AS loanDefaultedDate, l.loan_balance AS loanBal, c.full_name AS customerName, c.primary_mobile AS phoneNumber, c.national_id AS nationalId, c.branch, ls.name As loanStatus, u.email AS agentEmail")

		query = query.Joins("LEFT JOIN o_customers c ON l.customer_id = c.uid")
		query = query.Joins("LEFT JOIN o_users u ON u.uid = l.current_agent")
		query = query.Joins("LEFT JOIN o_loan_statuses ls ON ls.uid = l.status")
	}

	return query
}
