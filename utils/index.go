package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"super-lender/inits"
	"super-lender/models"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func QueryParamToIntWithDefault(c *gin.Context, param string, defaultValue int) int {
	intValue := TrimInt(c.Query(param))
	if intValue == 0 {
		return defaultValue
	}
	return intValue
}

func PathParamToIntWithDefault(c *gin.Context, param string, defaultValue int) int {
	intValue := TrimInt(c.Param(param))
	if intValue == 0 {
		return defaultValue
	}
	return intValue

}

func PathParamToStringWithDefault(c *gin.Context, param string, defaultValue string) string {
	paramValue := TrimString(c.Param(param))
	if paramValue == "" {
		return defaultValue
	}
	return paramValue
}

func QueryParamToStringWithDefault(c *gin.Context, param string, defaultValue string) string {
	paramValue := TrimString(c.Query(param))
	if paramValue == "" {
		return defaultValue
	}
	return paramValue
}

func GetDBConn(c *gin.Context) *gorm.DB {
	dbType, _ := c.Get("db")
	archive := os.Getenv("ARCHIVE")
	archiveVal, _ := strconv.Atoi(archive)

	if dbType == "archive" && archiveVal == 1 {
		return inits.ArchiveDB
	}

	return inits.CurrentDB
}

func ZeroToOne(val int) int {
	if val == 1 {
		return 0
	}
	return 1
}

func ToggleIcon(val int) string {
	if val == 1 {
		return "<i class=\"fa fa-times text-red\"></i>"
	}
	return "<i class=\"fa fa-check text-green\"></i>"
}

// TrimString trims leading and trailing spaces from a string
func TrimString(s string) string {
	return strings.TrimSpace(s)
}

// TrimInt trims leading and trailing spaces from an integer represented as a string
func TrimInt(s string) int {
	trimmed := strings.TrimSpace(s)
	intValue, err := strconv.Atoi(trimmed)
	if err != nil {
		// Handle error if conversion fails, e.g., return default value or panic
		return 0 // Default value
	}
	return intValue
}

// IsAvailable checks if the input is available (not empty or nil) and returns true, otherwise false
func IsAvailable(input interface{}) bool {
	if input == nil {
		return false
	}

	// Check for empty string
	if str, ok := input.(string); ok && str == "" {
		return false
	}

	// Add additional checks for other types if needed

	return true
}

func IsValidEmail(email string) bool {
	// Regular expression for email validation
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	match, _ := regexp.MatchString(regex, email)
	return match
}

func MakePhoneValid(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "+", "")

	// Read the country code from environment variable
	ccStr := os.Getenv("COUNTRY_CODE")
	cc := 0
	if ccStr != "" {
		fmt.Sscanf(ccStr, "%d", &cc)
	}

	if cc <= 0 {
		cc = 254
	}

	if len(phone) == 12 && phone[:3] == fmt.Sprintf("%d", cc) {
		return phone
	} else {
		if strings.HasPrefix(phone, "0") {
			phone = strings.TrimLeft(phone, "0")
			vphone := fmt.Sprintf("%d%s", cc, phone)
			return vphone
		} else {
			return fmt.Sprintf("%d%s", cc, phone)
		}
	}
}

func IsPhoneValid(phone string) bool {
	// Read the country code from environment variable
	ccStr := os.Getenv("COUNTRY_CODE")

	fmt.Println("Country code:", ccStr)
	cc := 0
	if ccStr != "" {
		fmt.Sscanf(ccStr, "%d", &cc)
	}

	if len(phone) == 12 && phone[:3] == fmt.Sprintf("%d", cc) {
		return true
	} else {
		return false
	}
}

// function to get UserGroup name from db based on the user group id
func GetUserGroupName(userGroupID int, userGroupStatus int) string {
	var userGroup models.OUserGroup
	result := inits.CurrentDB.Where("uid = ? AND status = ?", userGroupID, userGroupStatus).First(&userGroup)
	if result.Error != nil {
		return ""
	}
	return userGroup.Name
}

// GenerateRandomString generates a random string of given length
func GenerateRandomString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	const characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charactersLength := len(characters)
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = characters[rng.Intn(charactersLength)]
	}
	return string(randomString)
}

func GenerateRandomNumber(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	const characters = "0123456789"
	charactersLength := len(characters)
	randomNumber := make([]byte, length)

	for i := 0; i < length; i++ {
		randomNumber[i] = characters[rng.Intn(charactersLength)]
	}

	return string(randomNumber)
}

func CrazyString(length int) string {
	characters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#%^*()_+-~{}[];:|.<>"
	charactersLength := len(characters)
	randomString := make([]byte, length)

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for i := 0; i < length; i++ {
		randomString[i] = characters[rng.Intn(charactersLength)]
	}
	return string(randomString)
}

func CompanySettings() map[string]string {
	settings := make(map[string]string)
	var companySetting models.PlatformSetting

	// Assuming inits.CurrentDB is your Gorm database instance
	result := inits.CurrentDB.First(&companySetting, 1)
	if result.Error != nil {
		fmt.Println("Error fetching company settings:", result.Error)
		return settings
	}

	// Populate the settings map
	settings["Name"] = companySetting.Name
	settings["CompanyID"] = strconv.Itoa(companySetting.CompanyID)
	settings["Logo"] = companySetting.Logo
	settings["Icon"] = companySetting.Icon
	settings["Link"] = companySetting.Link

	return settings
}

func StoreEvent(tbl string, fld int, eventDetails string, eventBy int) {
	event := models.OEvent{
		Tbl:          tbl,
		Fld:          fld,
		EventDetails: eventDetails,
		EventBy:      eventBy,
	}

	result := inits.CurrentDB.Create(&event)
	if result.Error != nil {
		fmt.Println("Error storing event:", result.Error)
	}
}

func GetPermission(userId int, tbl string, rec int, act string) bool {
	db := inits.CurrentDB

	// Fetch user information
	var user models.OUser
	if err := db.First(&user, userId).Error; err != nil {
		fmt.Println("Error fetching user:", err)
		return false
	}

	// Check if the user is an admin
	if user.UserGroup == 1 {
		return true
	}

	// Query for permissions
	var permissions []models.OPermission
	query := db.Table("o_permissions")
	query = query.Where("tbl = ? AND rec = ? AND ? = ?", tbl, rec, act, 1)
	query = query.Where("group_id = ? OR user_id = ?", user.UserGroup, user.UID)
	if err := query.Find(&permissions).Error; err != nil {
		fmt.Println("Error fetching permissions:", err)
		return false
	}

	// Return true if any permission is found
	return len(permissions) > 0
}

func GeneratePlaceholders(length int) string {
	var placeholders []string
	for i := 0; i < length; i++ {
		placeholders = append(placeholders, "?")
	}
	return strings.Join(placeholders, ", ")
}

func IntSliceContains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func GetBranches(c *gin.Context, user models.OUser, readAll bool) []int {
	var branches []int
	if readAll {
		branches = append(branches, 0)
	} else {
		userBranch := user.Branch
		db := GetDBConn(c)
		db.Table("o_staff_branches").Where("agent = ? AND status = 1", user.UID).Select("branch").Scan(&branches)
		if userBranch > 0 {
			found := false
			for _, branch := range branches {
				if branch == userBranch {
					found = true
					break
				}
			}
			if !found {
				branches = append(branches, userBranch)
			}
		}
	}
	return branches
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lt":
		return "Should be less than " + fe.Param()
	case "lte":
		return "Should be less than " + fe.Param()
	case "gt":
		return "Should be greater than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "email":
		return "Invalid email address"
	case "min":
		return "Minimum length is " + fe.Param()
	case "max":
		return "Maximum length is " + fe.Param()
	case "alpha":
		return "Should contain only alphabets"
	case "alphanum":
		return "Should contain only alphabets and numbers"
	case "numeric":
		return "Should contain only numbers"
	case "eq":
		return "Should be equal to " + fe.Param()
	case "oneof":
		return "Invalid value. Allowed values are " + fe.Param()
	case "url":
		return "Invalid URL"
	case "datetime":
		return "Invalid date format"
	}
	return "Unknown error"
}

func IsValidFullName(fullName string) bool {
	// Regular expression pattern to match names with a single space
	pattern := `^[a-zA-Z]+ [a-zA-Z]+$`
	matched, _ := regexp.MatchString(pattern, fullName)
	return matched
}

func ValidateFullName(fl validator.FieldLevel) bool {
	fullName := fl.Field().String()
	// Regular expression pattern to match names with a single space
	pattern := `^[a-zA-Z]+ [a-zA-Z]+$`
	matched, _ := regexp.MatchString(pattern, fullName)
	return matched
}

func LogEvent(tbl string, fld int, eventDetails string, eventBy int) {
	event := models.OEvent{
		Tbl:          tbl,
		Fld:          fld,
		EventDetails: eventDetails,
		EventBy:      eventBy,
	}

	result := inits.CurrentDB.Create(&event)
	if result.Error != nil {
		fmt.Println("Error storing event:", result.Error)
	}
}

func CreateChangesLog(action, secAffectedTable, primaryAffectedEntityName string, primaryAffectedEntityID, secAffectedEntityID int, originalData, newData interface{}, user models.OUser, ignoredFields []string) {
	// Identify modified fields for update action
	var modifiedFields map[string]map[string]interface{}
	var eventDetails string

	switch action {
	case "update":
		modifiedFields = IdentifyModifiedFields(originalData, newData, ignoredFields)
		eventDetails = GenerateEventDetails("update", primaryAffectedEntityName, primaryAffectedEntityID, secAffectedEntityID, modifiedFields, user)
	case "delete":
		// Handle delete action, if needed
		eventDetails = fmt.Sprintf("Deletion of %s(%v) triggered by %s(%s)(UID: %d)", primaryAffectedEntityName, primaryAffectedEntityID, user.Name, user.Email, user.UID)
	default:
		fmt.Println("Unsupported action")
		return
	}

	LogEvent(secAffectedTable, secAffectedEntityID, eventDetails, user.UID)
}

func GenerateEventDetails(action string, primaryAffectedEntityName string, primaryAffectedEntityID, secAffectedEntityID int, modifiedFields map[string]map[string]interface{}, user models.OUser) string {
	if action == "update" && len(modifiedFields) > 0 {
		var logMessages []string
		for fieldName, values := range modifiedFields {
			oldVal := values["old"]
			newVal := values["new"]
			logMessages = append(logMessages, fmt.Sprintf("%s from %+v to %+v", fieldName, oldVal, newVal))
		}
		return fmt.Sprintf("%s(%v) %s triggered by [%s(%s)(%d)]. Changes: %s", primaryAffectedEntityName, primaryAffectedEntityID, action, user.Name, user.Email, user.UID, strings.Join(logMessages, ", "))
	}
	return fmt.Sprintf("%s(%v) %s triggered by [%s(%s)(%d)]. No values were modified", primaryAffectedEntityName, primaryAffectedEntityID, action, user.Name, user.Email, user.UID)
}

func AreEqual(a, b interface{}) bool {
	// Get the types of a and b
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)

	// If the types are different, they cannot be equal
	if typeA != typeB {
		return false
	}

	// Get the values of a and b
	valueA := reflect.ValueOf(a)
	valueB := reflect.ValueOf(b)

	// Iterate through the fields of the structs
	for i := 0; i < typeA.NumField(); i++ {
		// Get the fields of a and b
		fieldA := valueA.Field(i)
		fieldB := valueB.Field(i)

		// Compare the values of the fields
		if !reflect.DeepEqual(fieldA.Interface(), fieldB.Interface()) {
			return false
		}
	}

	// If all fields are equal, the structs are equal
	return true
}

// Function to identify modified fields and their old and new values
func IdentifyModifiedFields(original, updated interface{}, skipFields []string) map[string]map[string]interface{} {
	modifiedFields := make(map[string]map[string]interface{})

	// Map to store fields to skip for faster lookup
	skipMap := make(map[string]bool)
	for _, field := range skipFields {
		skipMap[field] = true
	}

	// Iterate through the fields of the structs
	for i := 0; i < reflect.TypeOf(original).NumField(); i++ {
		// Get the field name
		fieldName := reflect.TypeOf(original).Field(i).Name

		// Check if the field should be skipped
		if skipMap[fieldName] {
			continue
		}

		// Get the fields of original and updated customers
		fieldOriginal := reflect.ValueOf(original).Field(i)
		fieldUpdated := reflect.ValueOf(updated).Field(i)

		// Compare the values of the fields
		if !reflect.DeepEqual(fieldOriginal.Interface(), fieldUpdated.Interface()) {
			// Store the modified field and its old and new values
			modifiedFields[fieldName] = map[string]interface{}{
				"old": fieldOriginal.Interface(),
				"new": fieldUpdated.Interface(),
			}
		}
	}

	return modifiedFields
}

func CompareAndLogUpdates(original, updated interface{}) string {
	var buffer bytes.Buffer
	originalValue := reflect.ValueOf(original)
	updatedValue := reflect.ValueOf(updated)

	if originalValue.Kind() != updatedValue.Kind() {
		buffer.WriteString("Error: Original and updated values must be of the same type.\n")
		return buffer.String()
	}

	switch originalValue.Kind() {
	case reflect.Struct:
		for i := 0; i < originalValue.NumField(); i++ {
			originalFieldValue := originalValue.Field(i)
			updatedFieldValue := updatedValue.Field(i)

			if !reflect.DeepEqual(originalFieldValue.Interface(), updatedFieldValue.Interface()) {
				buffer.WriteString(fmt.Sprintf("%s has been updated from %v to %v\n", originalValue.Type().Field(i).Name, originalFieldValue.Interface(), updatedFieldValue.Interface()))
			}
		}

	case reflect.Map:
		for _, key := range originalValue.MapKeys() {
			originalMapValue := originalValue.MapIndex(key)
			updatedMapValue := updatedValue.MapIndex(key)

			if !reflect.DeepEqual(originalMapValue.Interface(), updatedMapValue.Interface()) {
				buffer.WriteString(fmt.Sprintf("Key %v has been updated from %v to %v\n", key.Interface(), originalMapValue.Interface(), updatedMapValue.Interface()))
			}
		}

	default:
		buffer.WriteString("Error: Unsupported data type.\n")
	}

	return buffer.String()
}

func CustomerStatusName(status int) string {
	switch status {
	case 0:
		return "Deleted"
	case 1:
		return "Active"
	case 2:
		return "Blocked"
	case 3:
		return "Lead"
	case 4:
		return "Draft"
	default:
		return "Unknown"
	}
}

func Sha256Hash(input string) string {
	hashed := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hashed[:])
}

func GetCountryCode() (string, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	// Get the value of COUNTRY_CODE from the environment
	countryCode := os.Getenv("COUNTRY_CODE")
	if countryCode == "" {
		return "", errors.New("COUNTRY_CODE not found in .env file")
	}

	return countryCode, nil
}

func InInterfaceArray(needle interface{}, haystack []interface{}) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func TruncateString(input string, length int) string {
	if len(input) <= length {
		return input
	}
	return input[:length]
}

func PascalCaseToSeparatedWords(input string) string {
	var words []string
	var currentWord string

	for i, char := range input {
		if i == 0 {
			currentWord += string(unicode.ToLower(char))
		} else {
			if unicode.IsUpper(char) {
				words = append(words, currentWord)
				currentWord = string(unicode.ToLower(char))
			} else {
				currentWord += string(char)
			}
		}
	}

	if len(currentWord) > 0 {
		words = append(words, currentWord)
	}

	return strings.Join(words, " ")
}
