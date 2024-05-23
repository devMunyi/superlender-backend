package controllers

import (
	"errors"
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

func GetCustomerContact(c *gin.Context) {
	// set result schema
	var customerContactResult schemas.GetCustomerContactResultSchema

	// Fetch query parameters from /customers/:uid
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
	query := db.Table("o_customer_contacts c")
	query = query.Select("c.uid, c.customer_id, c.contact_type, c.value, c.enc_phone, c.last_update")
	query = query.Where("c.uid = ?", uid)

	// execute query
	err := query.Scan(&customerContactResult).Error
	if err != nil {
		// send json response with error message
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// send json response
	c.JSON(200, customerContactResult)

}

func GetCustomerContacts(c *gin.Context) {

	// define a struct to hold the result
	var customerContactResult []schemas.GetCustomerContactsResultSchema

	// Fetch query parameters from /customers/:uid/contacts
	uid := utils.PathParamToIntWithDefault(c, "uid", 0)

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID
	readAll := utils.GetPermission(userId, "o_customers", 0, "read_")
	branches := utils.GetBranches(c, user, readAll)

	// Build query
	db := utils.GetDBConn(c)
	query := db.Table("o_customer_contacts c")
	query = query.Select("c.uid, c.contact_type, c.value, DATE_FORMAT(c.last_update, '%Y-%m-%d %H:%i:%s') AS last_update")

	// Apply filters
	query = query.Where("c.customer_id = ?", uid)
	if !readAll {
		query = query.Where("c.branch IN (?)", branches)
	}

	// Execute query
	err := query.Scan(&customerContactResult).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"contacts": customerContactResult,
	})
}

func CreateCustomerContact(c *gin.Context) {

	var customerContactInput models.OCustomerContacts
	var isEmail bool
	var isPhone bool
	var contactValueHash string

	if err := c.ShouldBindJSON(&customerContactInput); err != nil {
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

	contactValue := strings.TrimSpace(customerContactInput.Value)
	customerContactInput.Value = contactValue

	// if contact type is 1 or 2, validate phone number
	contactType := customerContactInput.ContactType
	if contactType == int(models.AlternativePhone1) || contactType == int(models.AlternativePhone2) {
		// make the value a valid phone number then check its validity
		contactValue = utils.MakePhoneValid(customerContactInput.Value)

		if !utils.IsPhoneValid(contactValue) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid phone number",
			})
			return

		}

		// at this point we know its a valid phone number
		isPhone = true

	}

	// if contact type is 3, validate email address
	if contactType == int(models.AlternativeEmail1) {
		// validate email address
		if !utils.IsValidEmail(contactValue) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid email address",
			})
			return
		}

		// at this point we know its a valid email address
		isEmail = true
	}

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID

	// check if user has necessary permissions permission to create customer contact
	createPermi := utils.GetPermission(userId, "o_customer_contacts", 0, "create_")
	if !createPermi {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  403,
			"message": "You don't have permission to create customer contact!",
		})
		return
	}

	// set db connection
	db := inits.CurrentDB

	// if isPhone is true, check if user with the same phone number exists in o_customers
	if isPhone {
		var existingCustomer models.OCustomer
		if err := db.Where("primary_mobile = ?", contactValue).First(&existingCustomer).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Phone number exists as primary mobile for a customer!",
			})
			return
		}

		// check from o_customer_contacts
		var existingCustomerContact models.OCustomerContacts
		if err := db.Where("value = ?", contactValue).First(&existingCustomerContact).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Phone number exists as alternative phone for a customer!",
			})
			return

		}

		// hash the phone number
		contactValueHash = utils.Sha256Hash(contactValue)

	}

	// if isEmail is true, check if user with the same email address exists in o_customers
	if isEmail {
		var existingCustomer models.OCustomer
		if err := db.Where("email_address = ?", contactValue).First(&existingCustomer).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Email address already exists!",
			})
			return
		}

		// check from o_customer_contacts
		var existingCustomerContact models.OCustomerContacts
		if err := db.Where("value = ?", contactValue).First(&existingCustomerContact).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Email address exists as alternative email for a customer!",
			})
			return

		}

		contactValueHash = ""
	}

	// Create customer contact
	customerContact := models.OCustomerContacts{
		CustomerID:  customerContactInput.CustomerID,
		ContactType: int(contactType), // Convert ContactType to int
		Value:       contactValue,
		EncPhone:    &contactValueHash, // Change type to *string
		LastUpdate:  time.Now().Local().Format("2006-01-02 15:04:05"),
	}

	if err := db.Create(&customerContact).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer contact created successfully",
		"contact": customerContact,
	})

}

func UpdateCustomerContact(c *gin.Context) {

	// implementation will be similar to that of creating a customer contact only that now were updating
	var customerContactInput models.OCustomerContacts
	var isEmail bool
	var isPhone bool
	if err := c.ShouldBindJSON(&customerContactInput); err != nil {
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

	contactValue := strings.TrimSpace(customerContactInput.Value)

	// if contact type is 1 or 2, validate phone number
	contactType := customerContactInput.ContactType
	customerContactInput.Value = contactValue
	if contactType == int(models.AlternativePhone1) || contactType == int(models.AlternativePhone2) {
		// make the value a valid phone number then check its validity
		contactValue = utils.MakePhoneValid(customerContactInput.Value)

		if !utils.IsPhoneValid(contactValue) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid phone number",
			})
			return

		}

		// at this point we know its a valid phone number
		isPhone = true
	}

	// if contact type is 3, validate email address
	if contactType == int(models.AlternativeEmail1) {
		// validate email address
		if !utils.IsValidEmail(contactValue) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid email address",
			})
			return
		}

		// at this point we know its a valid email address
		isEmail = true
	}

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID

	// check if user has necessary permissions permission to update customer contact
	updatePermi := utils.GetPermission(userId, "o_customer_contacts", 0, "update_")
	if !updatePermi {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  403,
			"message": "You don't have permission to update customer contact!",
		})
		return
	}

	// set db connection
	db := inits.CurrentDB

	// if isPhone is true, check if user with the same phone number exists in o_customers and not the current customer

	var contactValueHash string
	if isPhone {
		var existingCustomer models.OCustomer
		if err := db.Where("primary_mobile = ? AND uid != ?", contactValue, customerContactInput.CustomerID).First(&existingCustomer).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Phone number exists as primary mobile for another customer!",
			})
			return
		}

		// check from o_customer_contacts
		var existingCustomerContact models.OCustomerContacts
		if err := db.Where("value = ? AND uid != ?", contactValue, customerContactInput.UID).First(&existingCustomerContact).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Phone number exists as alternative phone for another customer!",
			})
			return

		}

		// hash the phone number
		contactValueHash = utils.Sha256Hash(contactValue)
	}

	// if isEmail is true, check if user with the same email address exists in o_customers and not the current customer
	if isEmail {
		var existingCustomer models.OCustomer
		if err := db.Where("email_address = ? AND uid != ?", contactValue, customerContactInput.CustomerID).First(&existingCustomer).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Email address already taken!",
			})
			return
		}

		// check from o_customer_contacts
		var existingCustomerContact models.OCustomerContacts
		if err := db.Where("value = ? AND uid != ?", contactValue, customerContactInput.UID).First(&existingCustomerContact).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Email address exists as alternative email for another customer!",
			})
			return

		}

		contactValueHash = ""
	}

	// Update customer contact

	customerContactInput.EncPhone = &contactValueHash
	customerContactInput.LastUpdate = time.Now().Local().Format("2006-01-02 15:04:05")

	// get original contact details before inserting new details
	var originalContact models.OCustomerContacts
	if err := db.Model(&models.OCustomerContacts{}).Where("uid = ?", customerContactInput.UID).First(&originalContact).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return

	}

	//  now update the contact details
	if err := db.Model(&models.OCustomerContacts{}).Where("uid = ?", customerContactInput.UID).Updates(&customerContactInput).Error; err != nil {
		// handle not found error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Customer contact not found",
			})
			return

		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	// log changes
	utils.CreateChangesLog("update", "o_customers", "Contact", customerContactInput.UID, customerContactInput.CustomerID, originalContact, customerContactInput, user, []string{"UID", "EncPhone", "LastUpdate"})

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer contact updated successfully",
		"contact": originalContact.UID,
	})

}
