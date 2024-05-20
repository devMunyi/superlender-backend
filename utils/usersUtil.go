package utils

import "gorm.io/gorm"

func FindManyUsersQueryBuilder(db *gorm.DB, userGroup, branch, status int, searchTerm, queryType string) *gorm.DB {
	query := db.Table("o_users u")
	if queryType == "count" {
		query = query.Select("u.uid")
	} else {
		query = query.Joins("LEFT JOIN o_user_groups ug ON u.user_group = ug.uid")
		query = query.Joins("LEFT JOIN o_staff_statuses ss ON u.status = ss.uid")
		query = query.Select("u.uid, u.name, u.email, DATE_FORMAT(u.join_date, '%Y-%m-%d %H:%i:%s') AS join_date, ug.name AS user_group, ss.name AS status")
	}

	if userGroup > 0 {
		query = query.Where("u.user_group = ?", userGroup)
	}
	if branch > 0 {
		query = query.Where("u.branch = ?", branch)
	}
	if status > 0 {
		query = query.Where("u.status = ?", status)
	}
	if searchTerm != "" {
		query = query.Where("u.name LIKE ? OR u.email LIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%")
	}

	return query
}
