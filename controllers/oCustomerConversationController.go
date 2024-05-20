package controllers

import (
	"super-lender/inits"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCustomerConversations(c *gin.Context) {
	var (
		pageNo     int
		rpp        int
		orderby    string
		dir        string
		searchTerm string
		countLimit int
	)

	// Fetch query parameters
	pageNo, _ = strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	rpp, _ = strconv.Atoi(c.DefaultQuery("rpp", "10"))
	orderby = c.DefaultQuery("orderby", "uid")
	dir = c.DefaultQuery("dir", "DESC")
	searchTerm = c.DefaultQuery("searchTerm", "")
	countLimit, _ = strconv.Atoi(c.DefaultQuery("countLimit", "1000"))

	type ConversationResult struct {
		UID              int       `json:"uid"`
		Transcript       string    `json:"transcript"`
		ConversationDate time.Time `json:"conversation_date"`
		NextInteraction  time.Time `json:"next_interaction"`
	}

	// Fetch 10 rows ordered by UID column in descending order with an offset
	offset := (pageNo - 1) * rpp
	query := `
        SELECT cc.uid, c.full_name, c.branch, cc.transcript, cc.conversation_date, cc.next_interaction
        FROM o_customer_conversations cc
        INNER JOIN o_customers c ON c.uid = cc.customer_id
        WHERE c.full_name LIKE ?
        ORDER BY cc.` + orderby + ` ` + dir + `
        LIMIT ? OFFSET ?
    `

	// Execute raw SQL query to fetch data
	var customerConversations []ConversationResult
	result := inits.CurrentDB.Raw(query, "%"+searchTerm+"%", rpp, offset).Scan(&customerConversations)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	// // Fetch total count of rows with limit
	countQuery := `
	SELECT cc.uid
	FROM o_customer_conversations cc
	INNER JOIN o_customers c ON c.uid = cc.customer_id
	WHERE c.full_name LIKE ?
	LIMIT ?`

	var totalCount int64
	countResult := inits.CurrentDB.Raw(countQuery, "%"+searchTerm+"%", countLimit).Count(&totalCount)

	// countQuery := `
	// SELECT COUNT(cc.uid)
	// FROM o_customer_conversations cc
	// INNER JOIN o_customers c ON c.uid = cc.customer_id
	// WHERE c.full_name LIKE ?`

	// countResult := inits.CurrentDB.Raw(countQuery, "%"+searchTerm+"%").Scan(&totalCount)

	if countResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"count": totalCount, // 57984
		"data":  customerConversations,
	})
}
