package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"super-lender/inits"
	"super-lender/models"
	"super-lender/schemas"
	customTypes "super-lender/types"
	"super-lender/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func FindManyCustomers(c *gin.Context) {

	// set necessary variables
	var customerResult []schemas.GetCustomersResultSchema
	var customerUIDCountResultSet []schemas.UIDCountResultsSchema
	db := utils.GetDBConn(c)
	pageNo := utils.QueryParamToIntWithDefault(c, "pageNo", 1)
	pageSize := utils.QueryParamToIntWithDefault(c, "pageSize", 10)
	orderBy := utils.QueryParamToStringWithDefault(c, "orderBy", "uid")
	dir := utils.QueryParamToStringWithDefault(c, "dir", "DESC")
	searchTerm := utils.QueryParamToStringWithDefault(c, "searchTerm", "")
	countLimit := utils.QueryParamToIntWithDefault(c, "countLimit", 0)
	branch := utils.QueryParamToIntWithDefault(c, "branch", 0)
	agent := utils.QueryParamToIntWithDefault(c, "agent", 0)
	status := utils.QueryParamToIntWithDefault(c, "status", 0)

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID
	readAll := utils.GetPermission(userId, "o_customers", 0, "read_")
	branches := utils.GetBranches(c, user, readAll)

	// Build select query
	selectQuery := utils.FindManyCustomersQueryBuilder(db, branch, agent, status, branches, readAll, searchTerm, "select")

	// Apply order and pagination
	selectQuery = selectQuery.Order("c." + orderBy + " " + dir)
	selectQuery = selectQuery.Limit(pageSize).Offset((pageNo - 1) * pageSize)

	// Execute selectQuery
	err := selectQuery.Scan(&customerResult).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	// count query
	countQuery := utils.FindManyCustomersQueryBuilder(db, branch, agent, status, branches, readAll, searchTerm, "count")
	var count int64
	if countLimit > 0 {
		/// ===== option 1: proves to work well when dealing large datasets
		countQuery = countQuery.Limit(countLimit)
		err = countQuery.Scan(&customerUIDCountResultSet).Error
		if err == nil {
			count = int64(len(customerUIDCountResultSet))
		}
	} else {
		// // ===== option 2: proves to work well when dealing small to medium datasets
		err = countQuery.Count(&count).Error
	}

	if err != nil {
		// send json response with error message
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"customers": customerResult,
		"count":     count,
	})
}

func FindCustomerById(c *gin.Context) {

	var customerResult schemas.GetCustomerResultSchema

	// Fetch query parameters from /customers/:uid
	uid := utils.PathParamToIntWithDefault(c, "uid", 0)

	primaryMobile := utils.QueryParamToStringWithDefault(c, "primary_mobile", "")

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID
	readAll := utils.GetPermission(userId, "o_customers", 0, "read_")
	branches := utils.GetBranches(c, user, readAll)

	// Build query
	db := utils.GetDBConn(c)
	query := db.Table("o_customers c")
	query = query.Select("c.uid, c.passport_photo, c.full_name, c.gender, c.dob, c.national_id, u.name AS added_by, c.loan_limit, c.email_address, c.primary_mobile, t.name AS phone_number_provider, b.name AS branch, c.physical_address, c.total_loans, p.name AS product, c.geolocation AS location_map, c.added_date, bd.icon AS badge_icon, bd.title AS badge_title, bd.description AS badge_description,  cs.name AS status")
	query = query.Joins("LEFT JOIN o_telecomms t ON c.phone_number_provider = t.uid")
	query = query.Joins("LEFT JOIN o_users u ON c.added_by = u.uid")
	query = query.Joins("LEFT JOIN o_branches b ON c.branch = b.uid")
	query = query.Joins("LEFT JOIN o_loan_products p ON c.primary_product = p.uid")
	query = query.Joins("LEFT JOIN o_badges bd ON c.badge_id = bd.uid")
	query = query.Joins("LEFT JOIN o_customer_statuses cs ON c.status = cs.code")

	// Apply filters
	if uid != 0 {
		query = query.Where("c.uid = ?", uid)
	} else if primaryMobile != "" {
		query = query.Where("c.primary_mobile = ?", primaryMobile)
	} else {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}

	if !readAll {
		query = query.Where("c.branch IN (?)", branches)
	}

	// Execute query
	err := query.Scan(&customerResult).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	// if uid == 0 return 404
	if customerResult.UID == 0 {
		c.JSON(404, gin.H{
			"message": "Customer not found",
		})
		return
	}

	customerResult.CurrentLO = customerResult.AddedBy
	customerResult.CurrentCO = customerResult.AddedBy

	// parse dob into date formatted as yyyy-mm-dd
	customerResult.Dob = utils.DateFormatter(customerResult.Dob)

	// format added_date into yyyy-mm-dd hh:mm:ss and should be in Nairobi TZ want 2024-04-01 10:48:08 not 2024-04-01T10:48:08
	customerResult.AddedDate = utils.DatetimeFormatter(customerResult.AddedDate)

	c.JSON(200, gin.H{
		"customer": customerResult,
	})
}

func CreateCustomer(c *gin.Context) {

	var createCustomerInput models.OCustomer

	if err := c.ShouldBindJSON(&createCustomerInput); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]customTypes.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = customTypes.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get country code
	countryCode, err := utils.GetCountryCode()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error: " + err.Error(),
		})
		return
	}

	if countryCode == "254" {
		// will use default of 1 for safaricom
	} else {
		// check if nonempty phone number provider is provided
		phoneNumberProvider := createCustomerInput.PhoneNumberProvider
		if phoneNumberProvider == 0 || phoneNumberProvider == 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Phone number provider is required",
			})
			return

		}
	}

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID

	// check if user with the same primary mobile exists
	db := inits.CurrentDB
	var otherWithSimilarPrimaryMobile models.OCustomer
	primaryMobile := utils.MakePhoneValid(createCustomerInput.PrimaryMobile)
	createCustomerInput.PrimaryMobile = primaryMobile
	if err := db.Where("primary_mobile = ?", primaryMobile).First(&otherWithSimilarPrimaryMobile).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Customer with the same primary mobile already exists",
		})
		return

	}

	// national id is not empty
	nationalId := strings.TrimSpace(createCustomerInput.NationalID)
	createCustomerInput.NationalID = nationalId
	if nationalId != "" {
		// check if user with the same national id exists
		var otherWithSimilarNationalId models.OCustomer
		if err := db.Where("national_id = ?", nationalId).First(&otherWithSimilarNationalId).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Customer with the same national id already exists",
			})
			return
		}

	}

	// if email address is not empty
	emailAddress := strings.TrimSpace(createCustomerInput.EmailAddress)
	createCustomerInput.EmailAddress = emailAddress
	if emailAddress != "" {
		// check if user with the same email address exists
		var otherWithSimilarEmail models.OCustomer
		if err := db.Where("email_address = ?", emailAddress).First(&otherWithSimilarEmail).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Customer with the same email address already exists",
			})
			return
		}

	}

	createCustomerInput.FullName = strings.TrimSpace(createCustomerInput.FullName)
	createCustomerInput.PhysicalAddress = strings.TrimSpace(createCustomerInput.PhysicalAddress)
	createCustomerInput.Geolocation = strings.TrimSpace(createCustomerInput.Geolocation)
	createCustomerInput.EncPhone = utils.Sha256Hash(primaryMobile)
	createCustomerInput.DOB = utils.DateFormatter(createCustomerInput.DOB)
	createCustomerInput.AddedBy = userId
	createCustomerInput.CurrentAgent = userId
	createCustomerInput.AddedDate = time.Now().Format("2006-01-02 15:04:05")

	if err := db.Create(&createCustomerInput).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Customer created successfully",
		"customer": createCustomerInput,
	})
}

func UpdateCustomer(c *gin.Context) {

	var updateCustomerInput models.OCustomer
	var existingCustomer models.OCustomer
	if err := c.ShouldBindJSON(&updateCustomerInput); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]customTypes.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = customTypes.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get country code
	countryCode, err := utils.GetCountryCode()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error: " + err.Error(),
		})
		return
	}

	if countryCode == "254" {
		fmt.Println("Country code is 254")
		// will use default of 1 for safaricom
	} else {
		// check if nonempty phone number provider is provided
		incomingPhoneNumberProvider := updateCustomerInput.PhoneNumberProvider

		// create an interface that holds 0, 1
		if incomingPhoneNumberProvider == 0 || incomingPhoneNumberProvider == 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Phone number provider is required",
			})
			return
		}

	}

	// set db connection
	db := inits.CurrentDB

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID

	//==== check if user has necessary permissions permission to update customer

	updatePermi := utils.GetPermission(userId, "o_customers", 0, "update_")
	if !updatePermi {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  403,
			"message": "You don't have permission to update customer!",
		})
		return
	}
	// retrieve original customer details for event logging
	if err := db.Model(&models.OCustomer{}).Where("uid = ?", updateCustomerInput.UID).First(&existingCustomer).Error; err != nil {
		// check if its a record not found error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status":  404,
				"message": "Customer not found",
			})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "Internal Server Error",
			})
			return
		}
	}

	// handle blocking && unblocking permission
	currentCustomerStatus := existingCustomer.Status
	incomingCustomerStatus := updateCustomerInput.Status
	if incomingCustomerStatus == models.BLOCKED {
		blockPermi := utils.GetPermission(userId, "o_customers", 0, "block_")
		if !blockPermi {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  403,
				"message": "You don't have permission to block customer!",
			})
			return
		}
	} else if incomingCustomerStatus == models.ACTIVE && currentCustomerStatus == models.BLOCKED {
		unblockPermi := utils.GetPermission(userId, "o_customers", 0, "unblock_")
		if !unblockPermi {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  403,
				"message": "You don't have permission to unblock customer!",
			})
			return
		}
	} else {
		// handle permission check for other statuses
		if incomingCustomerStatus != currentCustomerStatus {
			currentStatusName := utils.CustomerStatusName(int(currentCustomerStatus))
			updateStatusPermi := utils.GetPermission(userId, "o_customers", 0, "update_")
			if !updateStatusPermi {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"status":  403,
					"message": "You don't have permission to update customer who is " + currentStatusName,
				})
				return
			}
		}

	}

	// handle primary mobile update permission
	existingPrimaryMobile := existingCustomer.PrimaryMobile
	incomingPrimaryMobile := utils.MakePhoneValid(updateCustomerInput.PrimaryMobile)

	// check if phone is valid
	if !utils.IsPhoneValid(incomingPrimaryMobile) {

		fmt.Println("Invalid phone number", incomingPrimaryMobile)

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid phone number",
		})
		return
	} else {
		updateCustomerInput.PrimaryMobile = incomingPrimaryMobile
		updateCustomerInput.EncPhone = utils.Sha256Hash(incomingPrimaryMobile)
	}

	if incomingPrimaryMobile != existingPrimaryMobile {
		updatePrimaryPhonePermi := utils.GetPermission(userId, "o_customer_contacts", 0, "update_")

		if currentCustomerStatus == 1 && !updatePrimaryPhonePermi {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  403,
				"message": "You don't have permission to update customer phone number!",
			})
			return
		}

		// check if user has existing loan with the existingPrimaryMobile
		var existingLoan models.OLoan

		// fetch loan uid with the existingPrimaryMobile
		if err := db.Where("customer_id = ?", updateCustomerInput.UID).First(&existingLoan).Error; err == nil && existingLoan.UID != 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Customer has an existing loan",
			})
			return
		}

	}

	// ====== End of permission check

	// check if a user with the same primary mobile exists
	var otherWithSimilarPrimaryMobile models.OCustomer
	if err := db.Where("primary_mobile = ? AND uid != ?", incomingPrimaryMobile, updateCustomerInput.UID).First(&otherWithSimilarPrimaryMobile).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Customer with the same primary mobile already exists",
		})
		return
	}

	// check if a user with the same national id exists
	nationaId := strings.TrimSpace(updateCustomerInput.NationalID)
	updateCustomerInput.NationalID = nationaId
	if nationaId != "" {
		// check if user with the same national id exists
		var otherWithSimilarNationalId models.OCustomer
		if err := db.Where("national_id = ? AND uid != ?", nationaId, updateCustomerInput.UID).First(&otherWithSimilarNationalId).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Customer with the same national id already exists",
			})
			return
		}

	}

	// if email address is not empty
	emailAddress := strings.TrimSpace(updateCustomerInput.EmailAddress)
	updateCustomerInput.EmailAddress = emailAddress
	if emailAddress != "" {
		// check if user with the same email address exists
		var otherWithSimilarEmail models.OCustomer
		if err := db.Where("email_address = ? AND uid != ?", emailAddress, updateCustomerInput.UID).First(&otherWithSimilarEmail).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Customer with the same email address already exists",
			})
			return
		}
	}

	dobUpdate, err := utils.FormatDate(updateCustomerInput.DOB, 10)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid date format",
		})
		return
	} else {
		updateCustomerInput.DOB = dobUpdate
	}

	// store new customer details in the database
	if err := db.Model(&models.OCustomer{}).Where("uid = ?", updateCustomerInput.UID).Updates(&updateCustomerInput).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	// log customer update
	customerId := updateCustomerInput.UID
	existingCustomer.DOB = existingCustomer.DOB[:10]

	// log customer update
	utils.CreateChangesLog("update", "o_customers", "Customer", customerId, customerId, existingCustomer, updateCustomerInput, user, []string{"UID", "CurrentAgent", "AddedBy", "Geolocation", "AddedDate", "EncPhone"})

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Customer updated successfully",
		"data":    updateCustomerInput,
	})

}
