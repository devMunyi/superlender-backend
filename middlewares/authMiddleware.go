package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"super-lender/inits"
	"super-lender/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type UserResponse struct {
	Uid       int       `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	JoinDate  time.Time `json:"join_date"`
	UserGroup int       `json:"user_group"`
	Branch    int       `json:"branch"`
}

func RequireAuth(ctx *gin.Context) {
	// tokenString, err := ctx.Cookie("Authorization")
	// if err != nil {
	// 	ctx.JSON(401, gin.H{"error": "Failed to get Authorization cookie: " + err.Error()})
	// 	ctx.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	// fmt.Println("Authorization header: ", ctx.GetHeader("Authorization"))

	// extract token from header instead of cookie and remove Bearer prefix
	tokenString := ctx.GetHeader("Authorization")

	if tokenString == "" {
		ctx.JSON(401, gin.H{"error": "Authorization header is required"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return

	}

	// Remove Bearer prefix
	tokenString = tokenString[7:]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			errorMsg := fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])
			ctx.JSON(401, gin.H{"error": errorMsg})
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return nil, fmt.Errorf(errorMsg)
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		errorMsg := "Failed to parse token: " + err.Error()
		ctx.JSON(401, gin.H{"error": errorMsg})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["uid"], claims["exp"])
		scope := claims["scope"].([]interface{})
		db := scope[0].(string)

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.JSON(401, gin.H{"error": "Token expired"})
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.OUser
		// Attempt to convert claims["id"] directly to uint
		idFloat, ok := claims["uid"].(float64)
		if !ok {
			ctx.JSON(500, gin.H{"error": "Invalid ID type"})
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		id := int(idFloat) // Convert float64 to int

		if err := inits.CurrentDB.First(&user, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(401, gin.H{"error": "Unauthorized...."})
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.JSON(500, gin.H{"error": "Database error"})
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Set("user", user)
		ctx.Set("db", db)
	} else {
		ctx.JSON(401, gin.H{"error": "Invalid token claims"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	ctx.Next()
}
