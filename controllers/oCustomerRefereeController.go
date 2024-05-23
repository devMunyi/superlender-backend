package controllers

import (
	"errors"
	"net/http"
	"super-lender/models"
	"super-lender/schemas"
	customTypes "super-lender/types"
	"super-lender/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetCustomerReferee(c *gin.Context) {
	// set result schema
	var customerRefereeResult schemas.GetRefereeResultSchema

	// Fetch query parameters from /referees/:uid
	uid := utils.PathParamToIntWithDefault(c, "uid", 0)

	// if uid is 0, return error
	if uid == 0 {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID
	readPermi := utils.GetPermission(userId, "o_customers", 0, "read_")
	if !readPermi {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	// set db connection
	db := utils.GetDBConn(c)

	// build query
	query := db.Table("o_customer_referees r")
	query = query.Joins("LEFT JOIN o_customer_referee_relationships rr ON r.relationship = rr.uid")
	query = query.Select("r.uid, r.referee_name, r.customer_id, r.mobile_no, r.physical_address, rr.name AS relationship, r.added_date")
	query = query.Where("r.uid = ?", uid)

	// execute query
	err := query.Scan(&customerRefereeResult).Error
	if err != nil {

		// send json response with error message
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// send json response
	c.JSON(200, customerRefereeResult)

}

func GetCustomerReferees(c *gin.Context) {
	// set result schema
	var customerRefereesResult []schemas.GetRefereeResultSchema

	// Fetch query parameters from /customers/:uid/referees
	uid := utils.PathParamToIntWithDefault(c, "uid", 0)

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID
	readPermi := utils.GetPermission(userId, "o_customers", 0, "read_")
	if !readPermi {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	// set db connection
	db := utils.GetDBConn(c)

	// build query
	query := db.Table("o_customer_referees r")
	query = query.Joins("LEFT JOIN o_customer_referee_relationships rr ON r.relationship = rr.uid")
	query = query.Select("r.uid, r.referee_name, r.customer_id, r.mobile_no, r.physical_address, rr.name AS relationship, r.added_date")
	query = query.Where("r.customer_id = ?", uid)

	// execute query
	err := query.Scan(&customerRefereesResult).Error
	if err != nil {

		// send json response with error message
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// send json response
	c.JSON(200, customerRefereesResult)

}

func CreateCustomerReferee(c *gin.Context) {

	// set necessary variables
	var customerRefereeInput models.OCustomerReferee

	// bind request body to schema
	if err := c.ShouldBindJSON(&customerRefereeInput); err != nil {
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

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID
	writePermi := utils.GetPermission(userId, "o_customers", 0, "create_")
	if !writePermi {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	/// trim & sanitize inputs inputs
	customerRefereeInput.MobileNo = utils.MakePhoneValid(customerRefereeInput.MobileNo)

	// check if mobile number is valid
	if !utils.IsPhoneValid(customerRefereeInput.MobileNo) {
		c.JSON(400, gin.H{"error": "Invalid mobile number"})
		return
	}

	customerRefereeInput.PhysicalAddress = utils.TrimString(customerRefereeInput.PhysicalAddress)
	customerRefereeInput.RefereeName = utils.TrimString(customerRefereeInput.RefereeName)
	customerRefereeInput.EmailAddress = utils.TrimString(customerRefereeInput.EmailAddress)
	customerRefereeInput.AddedDate = time.Now().Local()

	// set db connection
	db := utils.GetDBConn(c)

	// check if referee with same customer_id, mobile_no, id_no, email_address and referee_name exists
	var refereeCount int64
	db.Table("o_customer_referees").Where("customer_id = ? AND mobile_no = ? AND referee_name = ? AND email_address = ?", customerRefereeInput.CustomerId, customerRefereeInput.MobileNo, customerRefereeInput.RefereeName, customerRefereeInput.EmailAddress).Count(&refereeCount)
	if refereeCount > 0 {
		c.JSON(400, gin.H{"error": "Referee already exists"})
		return

	}

	// execute query
	err := db.Create(&customerRefereeInput).Error
	if err != nil {
		// send json response with error message
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// send json response
	c.JSON(200, gin.H{"message": "Referee created successfully", "uid": customerRefereeInput.UID})
}

func UpdateCustomerReferee(c *gin.Context) {
	// set necessary variables
	var customerRefereeInput models.OCustomerReferee

	// bind request body to schema
	if err := c.ShouldBindJSON(&customerRefereeInput); err != nil {
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

	// check if UID is set
	if customerRefereeInput.UID == 0 {
		c.JSON(400, gin.H{"error": "UID is required"})
		return

	}

	///======= End of input validation

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID
	writePermi := utils.GetPermission(userId, "o_customers", 0, "update_")
	if !writePermi {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	/// trim & sanitize inputs
	customerRefereeInput.MobileNo = utils.MakePhoneValid(customerRefereeInput.MobileNo)

	// check if mobile number is valid
	if !utils.IsPhoneValid(customerRefereeInput.MobileNo) {
		c.JSON(400, gin.H{"error": "Invalid mobile number"})
		return
	}

	customerRefereeInput.PhysicalAddress = utils.TrimString(customerRefereeInput.PhysicalAddress)
	customerRefereeInput.RefereeName = utils.TrimString(customerRefereeInput.RefereeName)
	customerRefereeInput.EmailAddress = utils.TrimString(customerRefereeInput.EmailAddress)

	// set db connection
	db := utils.GetDBConn(c)

	// check if referee with same customer_id, mobile_no, id_no, email_address and referee_name exists
	var refereeCount int64
	refereeUID := customerRefereeInput.UID
	customerUID := customerRefereeInput.CustomerId
	db.Table("o_customer_referees").Where("customer_id = ? AND mobile_no = ? AND referee_name = ? AND email_address = ? AND uid != ?", customerUID, customerRefereeInput.MobileNo, customerRefereeInput.RefereeName, customerRefereeInput.EmailAddress, refereeUID).Count(&refereeCount)

	//  retrive stored referee to compare with the new referee for logging changes
	var existingRefereeDetails models.OCustomerReferee
	db.Table("o_customer_referees").Where("uid = ?", customerRefereeInput.UID).First(&existingRefereeDetails)

	// execute query
	err := db.Model(&customerRefereeInput).Where("uid = ?", customerRefereeInput.UID).Updates(&customerRefereeInput).Error

	if err != nil {
		// send json response with error message
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	utils.CreateChangesLog("update", "o_customers", "Referee", refereeUID, customerUID, existingRefereeDetails, customerRefereeInput, user, []string{"UID", "AddedDate", "IdNo", "Status"})

	// send json response
	c.JSON(200, gin.H{"message": "Referee updated successfully", "uid": refereeUID})

}
