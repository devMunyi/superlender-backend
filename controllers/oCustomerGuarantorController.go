package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

func GetCustomerGuarantor(c *gin.Context) {

	// set result schema
	var customerGuarantorResult schemas.GetCustomerGuarantorsResultSchema

	// Fetch query parameters from /customers/:uid
	uid := utils.PathParamToIntWithDefault(c, "uid", 0)

	// get user permission
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
	query := db.Table("o_customer_guarantors g")
	// join with o_customer_guarantor_relationships to get relationship name
	query = query.Joins("LEFT JOIN o_customer_guarantor_relationships gr ON gr.uid =g.relationship")
	query = query.Select("g.uid, g.guarantor_name, g.customer_id, g.mobile_no, g.national_id, g.physical_address, g.amount_guaranteed, g.added_date, gr.name AS relationship, g.status")

	// Apply filters
	query = query.Where("g.uid = ?", uid)

	// execute query
	err := query.Scan(&customerGuarantorResult).Error
	if err != nil {
		// send json response with error message
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// send json response
	c.JSON(200, customerGuarantorResult)
}

func CreateCustomerGuarantor(c *gin.Context) {

	var customerGuarantorInput schemas.CreateCustomerGuarantorSchema
	if err := c.ShouldBindJSON(&customerGuarantorInput); err != nil {
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

	// check if user has necessary permissions permission to create customer guarantor
	createPermi := utils.GetPermission(userId, "o_customers", 0, "update_")
	if !createPermi {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  403,
			"message": "You don't have permission to create customer guarantor!",
		})
		return
	}

	// set db connection
	db := inits.CurrentDB
	guarantorName := utils.TrimString(customerGuarantorInput.GuarantorName)
	mobileNo := utils.MakePhoneValid(customerGuarantorInput.MobileNo)
	nationalId := utils.TrimString(customerGuarantorInput.NationalId)
	physicalAddress := utils.TrimString(customerGuarantorInput.PhysicalAddress)

	// validate mobile number
	if !utils.IsPhoneValid(mobileNo) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid phone number",
		})
		return

	}

	// check existence of guarantor with same guarantor_name, national_id, mobile_no and customer_id
	var existingGuarantor models.OCustomerGuarantor
	if err := db.Where("guarantor_name = ? AND national_id = ? AND mobile_no = ? AND customer_id = ?", guarantorName, nationalId, mobileNo, customerGuarantorInput.CustomerId).First(&existingGuarantor).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Guarantor already exists",
		})
		return
	}

	// check uniqueness of national_id
	var existingNationalIDGuarantor models.OCustomerGuarantor
	if err := db.Where("national_id = ?", nationalId).First(&existingNationalIDGuarantor).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "national id already exists",
		})
		return
	}

	// check uniqueness of mobile_no
	var existingMobileNoGuarantor models.OCustomerGuarantor
	if err := db.Where("mobile_no = ?", mobileNo).First(&existingMobileNoGuarantor).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "mobile no already exists",
		})
		return
	}

	// Create customer guarantor
	customerGuarantor := models.OCustomerGuarantor{
		GuarantorName:    guarantorName,
		CustomerId:       customerGuarantorInput.CustomerId,
		MobileNo:         mobileNo,
		NationalId:       nationalId,
		PhysicalAddress:  physicalAddress,
		AmountGuaranteed: customerGuarantorInput.AmountGuaranteed,
		AddedDate:        time.Now().Format("2006-01-02 15:04:05"),
		Relationship:     customerGuarantorInput.Relationship,
	}

	if err := db.Create(&customerGuarantor).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Customer guarantor created successfully",
		"guarantor": customerGuarantor,
	})

}

func UpdateCustomerGuarantor(c *gin.Context) {

	var customerGuarantorInput schemas.UpdateCustomerGuarantorSchema
	if err := c.ShouldBindJSON(&customerGuarantorInput); err != nil {
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

	// Check if user has necessary permissions to update customer guarantor
	updatePermi := utils.GetPermission(userId, "o_customers", 0, "update_")
	if !updatePermi {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  403,
			"message": "You don't have permission to update customer guarantor!",
		})
		return
	}

	// Set db connection
	db := inits.CurrentDB

	// Trim and validate input data
	guarantorName := utils.TrimString(customerGuarantorInput.GuarantorName)
	mobileNo := utils.MakePhoneValid(customerGuarantorInput.MobileNo)
	nationalId := utils.TrimString(customerGuarantorInput.NationalId)
	physicalAddress := utils.TrimString(customerGuarantorInput.PhysicalAddress)
	relationship := utils.TrimInt(strconv.Itoa(customerGuarantorInput.Relationship))
	status := utils.TrimInt(strconv.Itoa(customerGuarantorInput.Status))

	if !utils.IsPhoneValid(mobileNo) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid phone number",
		})
		return
	}

	// Check if the guarantor already exists (excluding the current guarantor)
	if err := db.Where("guarantor_name = ? AND national_id = ? AND mobile_no = ? AND customer_id = ? AND uid != ?",
		guarantorName, nationalId, mobileNo, customerGuarantorInput.CustomerId, customerGuarantorInput.UID).
		First(&models.OCustomerGuarantor{}).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Guarantor already exists",
		})
		return
	}

	// Check if the national ID is unique (excluding the current guarantor)
	if err := db.Where("national_id = ? AND uid != ?", nationalId, customerGuarantorInput.UID).
		First(&models.OCustomerGuarantor{}).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "National ID already exists",
		})
		return
	}

	// Check if the mobile number is unique (excluding the current guarantor)
	if err := db.Where("mobile_no = ? AND uid != ?", mobileNo, customerGuarantorInput.UID).
		First(&models.OCustomerGuarantor{}).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Mobile number already exists",
		})
		return
	}

	// Get the original guarantor details
	var originalGuarantor schemas.UpdateCustomerGuarantorSchema
	if err := db.Model(&models.OCustomerGuarantor{}).Where("uid = ?", customerGuarantorInput.UID).
		Select("guarantor_name, customer_id, national_id, mobile_no, physical_address, amount_guaranteed, relationship, status").
		First(&originalGuarantor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Guarantor not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	// Update the guarantor details
	customerGuarantor := models.OCustomerGuarantor{
		UID:              customerGuarantorInput.UID,
		GuarantorName:    guarantorName,
		CustomerId:       customerGuarantorInput.CustomerId,
		NationalId:       nationalId,
		MobileNo:         mobileNo,
		PhysicalAddress:  physicalAddress,
		AmountGuaranteed: customerGuarantorInput.AmountGuaranteed,
		Relationship:     relationship,
		Status:           status,
	}

	if err := db.Model(&models.OCustomerGuarantor{}).Where("uid = ?", customerGuarantorInput.UID).
		Updates(&customerGuarantor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Guarantor not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	// Identify modified fields
	modifiedFields := utils.IdentifyModifiedFields(originalGuarantor, customerGuarantor, []string{"UID", "AddedDate"})

	// Log the update event if there were any changes
	var eventDetails string
	if len(modifiedFields) > 0 {
		var logMessages []string
		for fieldName, values := range modifiedFields {
			oldVal := values["old"]
			newVal := values["new"]
			logMessages = append(logMessages, fmt.Sprintf("%s changed from %+v to %+v", fieldName, oldVal, newVal))
		}
		eventDetails = fmt.Sprintf("Update triggered by %s [UID: %d]. Changes: %s", user.Name, userId, strings.Join(logMessages, ", "))

	} else {
		eventDetails = fmt.Sprintf("Update triggered by %s [UID: %d]. No values were modified", user.Name, userId)
		fmt.Println("No fields were modified")
	}

	utils.LogEvent("o_customers", customerGuarantorInput.CustomerId, eventDetails, userId)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer guarantor updated successfully",
		"data": gin.H{
			"guarantorName":    customerGuarantor.GuarantorName,
			"customerId":       customerGuarantor.CustomerId,
			"nationalId":       customerGuarantor.NationalId,
			"mobileNo":         customerGuarantor.MobileNo,
			"physicalAddress":  customerGuarantor.PhysicalAddress,
			"amountGuaranteed": customerGuarantor.AmountGuaranteed,
			"relationship":     customerGuarantor.Relationship,
			"status":           customerGuarantor.Status,
		},
	})
}

func GetCustomerGuarantors(c *gin.Context) {
	// define a struct to hold the result
	var customerGuarantorResult []schemas.GetCustomerGuarantorsResultSchema

	// Fetch query parameters from /customers/:uid/guarantors
	uid := utils.PathParamToIntWithDefault(c, "uid", 0)

	// Get user and permissions
	user := c.MustGet("user").(models.OUser)
	userId := user.UID
	readPermi := utils.GetPermission(userId, "o_customers", 0, "read_")
	if !readPermi {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  403,
			"message": "You don't have permission to read customer guarantors!",
		})
		return
	}

	// set db connection
	db := utils.GetDBConn(c)

	// build query
	query := db.Table("o_customer_guarantors g")
	// join with o_customer_guarantor_relationships to get relationship name
	query = query.Joins("LEFT JOIN o_customer_guarantor_relationships gr ON gr.uid =g.relationship")
	query = query.Select("g.uid, g.guarantor_name, g.customer_id, g.mobile_no, g.national_id, g.physical_address, g.amount_guaranteed, g.added_date, gr.name AS relationship, g.status")

	// Apply filters
	query = query.Where("g.customer_id = ?", uid)

	// Execute query
	err := query.Scan(&customerGuarantorResult).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	// Apply the DatetimeFormatter to each element in customerGuarantorResult
	// for i := range customerGuarantorResult {
	// 	customerGuarantorResult[i].AddedDate = utils.DatetimeFormatter(customerGuarantorResult[i].AddedDate)
	// }

	c.JSON(200, gin.H{
		"guarantors": customerGuarantorResult,
	})
}
