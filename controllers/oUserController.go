package controllers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
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
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

func Signup(c *gin.Context) {
	// set necessary variables
	var userSignupInput models.OUser

	// bind request body to schema
	if err := c.ShouldBindJSON(&userSignupInput); err != nil {
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

	///// === begin validations

	// email validation
	email := strings.TrimSpace(userSignupInput.Email)
	userSignupInput.Email = email
	if email == "" {
		c.JSON(400, gin.H{"error": "Email is required"})
		return
	}

	// check if email is valid
	if !utils.IsValidEmail(email) {
		c.JSON(400, gin.H{"error": "Invalid email"})
		return
	}

	// check if user with email exists
	var user models.OUser
	result := inits.CurrentDB.Where("email = ?", email).First(&user)
	if result.Error == nil {
		c.JSON(400, gin.H{"error": "User with that email already exists"})
		return

	}
	///==== end of email validation

	///==== phone validation
	phone := utils.MakePhoneValid(userSignupInput.Phone)
	userSignupInput.Phone = phone
	// check if phone is valid
	if !utils.IsPhoneValid(phone) {
		c.JSON(400, gin.H{"error": "Invalid phone number"})
		return
	}

	// check if user with phone exists
	result = inits.CurrentDB.Where("phone = ?", phone).First(&user)
	if result.Error == nil {
		c.JSON(400, gin.H{"error": "User with that phone already exists"})
		return
	}

	//// ==== end of phone validation

	/// ==== national id validation
	nationalId := strings.TrimSpace(userSignupInput.NationalID)
	userSignupInput.NationalID = nationalId
	if nationalId == "" {
		c.JSON(400, gin.H{"error": "National ID is required"})
		return
	}

	// check if user with national id exists
	result = inits.CurrentDB.Where("national_id = ?", nationalId).First(&user)
	if result.Error == nil {
		c.JSON(400, gin.H{"error": "User with that national id already exists"})
		return
	}

	/// ==== end of national id validation

	//// ==== name validation
	name := strings.TrimSpace(userSignupInput.Name)
	userSignupInput.Name = name
	if name == "" {
		c.JSON(400, gin.H{"error": "Name is required"})
		return
	}

	/// ==== end of name validation

	/// ==== password validation
	password := strings.TrimSpace(userSignupInput.Pass1)
	userSignupInput.Pass1 = password
	// validate password length
	if len(password) < 6 {
		c.JSON(400, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	//// ==== end of password validation

	///==== validate user group

	if userSignupInput.UserGroup < 1 {
		c.JSON(400, gin.H{"error": "User group is required"})
		return

	}

	////=== end of user group validation

	///// ==== hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}

	userSignupInput.Pass1 = string(hash)
	userSignupInput.JoinDate = time.Now().Local()
	companyId, err := strconv.Atoi(os.Getenv("COMPANY_ID"))
	if err != nil {
		// provide a default company id as 1
		companyId = 1
	}
	userSignupInput.Company = companyId

	result = inits.CurrentDB.Create(&userSignupInput)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"data": "User created successfully", "uid": userSignupInput.UID})
}

func Login(c *gin.Context) {

	// set necessary variables
	var userLoginInput schemas.UserLoginSchema
	var user models.OUser

	// bind request body to schema
	if err := c.ShouldBindJSON(&userLoginInput); err != nil {
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

	/// === validate username
	username := strings.TrimSpace(userLoginInput.EmailOrPhone)
	if username == "" {
		c.JSON(400, gin.H{"error": "Username is required"})
		return
	}

	// check if username is email or phone
	isEmail := utils.IsValidEmail(username)
	if !isEmail {
		phone := utils.MakePhoneValid(username)
		if !utils.IsPhoneValid(phone) {
			c.JSON(400, gin.H{"error": "Invalid username or password"})
			return
		}

		result := inits.CurrentDB.Where("phone = ?", phone).First(&user)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Invalid username or password"})
			return
		}
	} else {
		result := inits.CurrentDB.Where("email = ?", username).First(&user)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Invalid username or password"})
			return
		}

	}

	/// === end of username validation

	/// === validate password

	password := strings.TrimSpace(userLoginInput.Password)
	if password == "" {
		c.JSON(400, gin.H{"error": "Password is required"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Pass1), []byte(password))

	if err != nil {
		fmt.Printf("Error comparing password: %v\n", err)
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	/// === end of password validation

	// === begin generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.UID,
		"scope": []string{"current", "read", "write"},
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"uid": user.UID,
	// 	"scope": map[string]interface{}{
	// 		"db": "current",
	// 		"op": []string{"read", "write"},
	// 	},
	// 	"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	// })

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	tokenResponse := schemas.UserTokenSchema{
		Token:     tokenString,
		ExpiresIn: 3600 * 24 * 30,
		TokenType: "Bearer",
		Scope:     []string{"current", "read", "write"},
	}

	// call GetUserGroupName function to get user group name
	userGroup := utils.GetUserGroupName(user.UserGroup, int(models.Active))

	userResponse := schemas.UserSchema{
		UID:       user.UID,
		Name:      user.Name,
		Email:     user.Email,
		JoinDate:  user.JoinDate.Local().Format("2006-01-02 15:04:05"),
		UserGroup: userGroup,
	}

	c.JSON(200, gin.H{"data": gin.H{"user": userResponse, "token": tokenResponse}})

	// // ctx.JSON(200, gin.H{"token": tokenString})
	// ctx.SetSameSite(http.SameSiteLaxMode)
	// ctx.SetCookie("Authorization", tokenString, 3600*24*30, "", "localhost", false, true)
	// ctx.SetCookie("db", "current", 3600*24*30, "", "localhost", false, true)
}

func FindManyUsers(c *gin.Context) {
	// set result schema
	var userResultSchema []schemas.UserSchema
	var userUIDCountResultSet []schemas.UIDCountResultsSchema

	// Fetch query parameters
	pageNo := utils.QueryParamToIntWithDefault(c, "pageNo", 1)
	pageSize := utils.QueryParamToIntWithDefault(c, "pageSize", 10)
	orderBy := utils.QueryParamToStringWithDefault(c, "orderBy", "uid")
	dir := utils.QueryParamToStringWithDefault(c, "dir", "DESC")
	searchTerm := utils.QueryParamToStringWithDefault(c, "searchTerm", "")
	countLimit := utils.QueryParamToIntWithDefault(c, "countLimit", 1000)
	userGroup := utils.QueryParamToIntWithDefault(c, "userGroup", 0)
	branch := utils.QueryParamToIntWithDefault(c, "branch", 0)
	status := utils.QueryParamToIntWithDefault(c, "status", 0)

	db := utils.GetDBConn(c)

	// build select query
	selectQuery := utils.FindManyUsersQueryBuilder(db, userGroup, branch, status, searchTerm, "select")
	selectQuery = selectQuery.Order(orderBy + " " + dir)
	selectQuery = selectQuery.Offset((pageNo - 1) * pageSize).Limit(pageSize)

	err := selectQuery.Scan(&userResultSchema).Error
	if err != nil {
		// send json response with error message
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// build count query
	countQuery := utils.FindManyUsersQueryBuilder(db, userGroup, branch, status, searchTerm, "count")
	var count int64
	if countLimit > 0 {
		/// ===== option 1: proves to work well when dealing large datasets
		countQuery = countQuery.Limit(countLimit)
		err = countQuery.Scan(&userUIDCountResultSet).Error
		count = int64(len(userUIDCountResultSet))
	} else {
		// // ===== option 2: proves to work well when dealing small to medium datasets
		err = countQuery.Count(&count).Error
	}

	if err != nil {
		// send json response with error message
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// send json response
	c.JSON(200, gin.H{"count": count, "data": userResultSchema})

}

// find one user uid from the database
func FindUserByID(c *gin.Context) {

	// set result schema
	var userResultSchema schemas.UserSchema

	// Fetch query parameters from /users/:uid
	uid := utils.PathParamToIntWithDefault(c, "uid", 0)

	// if uid is 0, return error
	if uid == 0 {
		c.JSON(400, gin.H{"error": "Invalid user id"})
		return
	}

	// set db connection
	db := utils.GetDBConn(c)

	// build query
	query := db.Table("o_users u")
	query = query.Joins("LEFT JOIN o_user_groups ug ON u.user_group = ug.uid")
	query = query.Joins("LEFT JOIN o_staff_statuses ss ON u.status = ss.uid")
	query = query.Select("u.uid, u.name, u.email, u.join_date, ug.name AS user_group, ss.name AS status")
	query = query.Where("u.uid = ?", uid)

	// execute query
	err := query.Scan(&userResultSchema).Error
	if err != nil {

		// send json response with error message
		c.JSON(400, gin.H{"error": err.Error()})
		return

	}

	// send json response
	userResultSchema.JoinDate = utils.DatetimeFormatter(userResultSchema.JoinDate)
	c.JSON(200, userResultSchema)

}

func ValidateUser(ctx *gin.Context) {
	// Retrieve the user from the context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert the user interface to User type
	userObj, ok := user.(models.OUser)
	if !ok {
		ctx.JSON(500, gin.H{"error": "Invalid user object type"})
		return
	}

	/// get user group name
	userGroup := utils.GetUserGroupName(userObj.UserGroup, int(models.Active))

	// Create a new UserResponse object with necessary fields
	userResponse := schemas.UserSchema{
		UID:       userObj.UID,
		Name:      userObj.Name,
		Email:     userObj.Email,
		JoinDate:  userObj.JoinDate.Local().Format("2006-01-02 15:04:05"),
		UserGroup: userGroup,
	}

	// Respond with user information
	ctx.JSON(200, gin.H{"data": "You are logged in!", "user": userResponse})
}

func Logout(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", "", -1, "", "localhost", false, true)
	ctx.JSON(200, gin.H{"data": "You are logged out!"})
}

func SwitchDB(ctx *gin.Context) {

	var Db struct {
		DbType string `json:"db_type"`
	}

	if err := ctx.ShouldBindJSON(&Db); err != nil {
		fmt.Println("Error binding JSON:", err)
		ctx.JSON(400, gin.H{"error": "Bad request"})
		return

	}

	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// db, exists := ctx.Get("db")
	// if !exists {
	// 	ctx.JSON(500, gin.H{"error": "A required parameter to switch is missing"})
	// 	return
	// }

	db := Db.DbType
	fmt.Println("DB Type:", db)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.(models.OUser).UID,
		"scope": []string{db, "read", "write"},
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		ctx.JSON(500, gin.H{"error": err})
		return
	}

	tokenResponse := schemas.UserTokenSchema{
		Token:     tokenString,
		ExpiresIn: 3600 * 24 * 30,
		TokenType: "Bearer",
		Scope:     []string{"current", "read", "write"},
	}

	// userResponse := UserResponse{
	// 	Uid:       user.(models.OUser).UID,
	// 	Name:      user.(models.OUser).Name,
	// 	Email:     user.(models.OUser).Email,
	// 	JoinDate:  user.(models.OUser).JoinDate,
	// 	UserGroup: utils.GetUserGroupName(user.(models.OUser).UserGroup, int(models.Active)),
	// }

	var message string = "You are now viewing live data"
	if db == "archive" {
		message = "You are now viewing archived data"
	}

	ctx.JSON(200, gin.H{"data": gin.H{"token": tokenResponse, message: message}})
}

func Login2(ctx *gin.Context) {
	var LoginBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&LoginBody); err != nil {
		fmt.Println("Error binding JSON:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	var user models.OUser
	result := inits.CurrentDB.Where("email = ?", LoginBody.Email).First(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	// Fetch the salt from the database for the user
	// Here you should have a function to fetch the salt based on the user's ID

	fmt.Println("UserID", user.UID)

	salt := fetchSaltFromDB(user.UID)

	fmt.Println("Salt:", salt)

	// Append salt to the input password
	fullPassword := salt + LoginBody.Password

	fmt.Println("Full Password:", fullPassword)

	// Hash the password with SHA256
	hashedPassword := sha256Hash(fullPassword)

	fmt.Println("Input User Hashed Password:", hashedPassword)
	fmt.Println("DB User Hashed Password:", user.Pass1)

	// Compare hashed passwords
	if user.Pass1 != hashedPassword {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.UID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Set JWT token as a cookie
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600*24*30, "", "localhost", false, true)
}

// Fetch the salt for the user from the database
func fetchSaltFromDB(userID int) string {
	var pass models.OPass

	// Execute the query using GORM's methods
	if err := inits.CurrentDB.Where("user = ?", userID).Select("pass").First(&pass).Error; err != nil {
		return ""
	}

	return pass.Pass
}

// SHA256 hashing function
func sha256Hash(input string) string {
	h := sha3.New256()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func ChangePassword(c *gin.Context) {
	var changePassword struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&changePassword); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	// // Retrieve the user from the context
	// user, exists := c.Get("user")
	// if !exists {
	// 	c.JSON(401, gin.H{"error": "Unauthorized"})
	// 	return
	// }

	// // Convert the user interface to User type
	// userObj, ok := user.(models.OUser)
	// if !ok {
	// 	c.JSON(500, gin.H{"error": "Invalid user object type"})
	// 	return
	// }

	// Fetch the user from the database
	var dbUser models.OUser
	result := inits.CurrentDB.Where("email = ?", changePassword.Email).First(&dbUser)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error fetching user"})
		return
	}

	// // Compare the old password
	// err := bcrypt.CompareHashAndPassword([]byte(dbUser.Pass1), []byte(changePassword.OldPassword))
	// if err != nil {
	// 	c.JSON(401, gin.H{"error": "Invalid password"})
	// 	return
	// }

	// Hash the new password
	hash, err := bcrypt.GenerateFromPassword([]byte(changePassword.NewPassword), 10)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}

	// Update the user's password
	result = inits.CurrentDB.Model(&dbUser).Update("pass1", string(hash))
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error updating password"})
		return
	}

	c.JSON(200, gin.H{"data": "Password updated successfully"})
}
