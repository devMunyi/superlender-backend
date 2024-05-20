package utils

import (
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func FindManyCustomersQueryBuilder(db *gorm.DB, branch, agent, status int, branches []int, readAll bool, searchTerm, queryType string) *gorm.DB {
	query := db.Table("o_customers c")

	if queryType == "count" {
		query = query.Select("c.uid")
	} else {
		query = query.Select("c.uid, c.passport_photo, c.full_name, u.name AS agent, c.email_address, c.primary_mobile, b.name AS branch, c.physical_address, cs.name AS status")
		query = query.Joins("LEFT JOIN o_users u ON c.current_agent = u.uid")
		query = query.Joins("LEFT JOIN o_branches b ON c.branch = b.uid")
		query = query.Joins("LEFT JOIN o_customer_statuses cs ON c.status = cs.code")
	}

	// Apply filters
	if branch != 0 {
		query = query.Where("c.branch = ?", branch)
	}
	if agent != 0 {
		query = query.Where("c.current_agent = ?", agent)
	}
	if status != 0 {
		query = query.Where("c.status = ?", status)
	}
	if !readAll {
		query = query.Where("c.branch IN (?)", branches)
	}

	// Apply search term
	if searchTerm != "" {
		if _, err := strconv.Atoi(searchTerm); err == nil {
			// searchTerm contains only digits
			///==search by phone number
			if len(searchTerm) == 12 && (strings.HasPrefix(searchTerm, "2547") || strings.HasPrefix(searchTerm, "2541")) {
				// searchTerm is a 12-digit number starting with 2547 or 2541
				query = query.Where("c.primary_mobile = ?", searchTerm)

			} else {
				// searchTerm by uid
				query = query.Where("c.uid = ?", searchTerm)
			}
		} else {
			isAlphabetic := regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(searchTerm)
			if isAlphabetic {
				// searchTerm contains only alphabets
				query = query.Where("c.full_name LIKE ?", "%"+searchTerm+"%")
			} else {
				// searchTerm is alphanumeric
				query = query.Where("c.uid LIKE ? OR c.primary_mobile LIKE ? OR c.full_name LIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")
			}
		}
	}
	return query
}
