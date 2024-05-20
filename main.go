package main

import (
	"super-lender/controllers"
	"super-lender/inits"
	"super-lender/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

func init() {
	inits.LoadEnv()
	inits.TZInit()
	inits.DBInit()
	inits.CurrentDB.Logger.LogMode(logger.Info)
	inits.ArchiveDB.Logger.LogMode(logger.Info)

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	_ = v.RegisterValidation("validFullName", validateFullName)
	// }
}

func main() {

	// Create a gin router
	r := gin.Default()

	////==== Begin general routes
	r.GET("/ping", controllers.HeathCheck)
	////==== End general routes

	////==== Begin users routes
	r.GET("/users", middlewares.RequireAuth, controllers.FindManyUsers)
	r.GET("/users/:uid", middlewares.RequireAuth, controllers.FindUserByID)
	r.POST("/users/signup", controllers.Signup)
	r.POST("/users/login", controllers.Login)
	r.GET("/users/auth", middlewares.RequireAuth, controllers.ValidateUser)
	r.GET("/users/logout", middlewares.RequireAuth, controllers.Logout)
	r.POST("/users/switch-db", middlewares.RequireAuth, controllers.SwitchDB)
	r.PUT("/users/change-password", controllers.ChangePassword)
	////==== End users routes

	////==== Begin customers routes
	r.POST("/customers", middlewares.RequireAuth, controllers.CreateCustomer)
	r.GET("/customers", middlewares.RequireAuth, controllers.FindManyCustomers)
	r.PUT("/customers", middlewares.RequireAuth, controllers.UpdateCustomer)
	r.GET("/customers/:uid", middlewares.RequireAuth, controllers.FindCustomerById)
	r.GET("/customers/:uid/contacts", middlewares.RequireAuth, controllers.GetCustomerContacts)
	r.GET("/customers/:uid/guarantors", middlewares.RequireAuth, controllers.GetCustomerGuarantors)
	r.GET("/customers/:uid/referees", middlewares.RequireAuth, controllers.GetCustomerReferees)
	////==== End customers routes

	////==== Begin contacts routes
	r.POST("/contacts", middlewares.RequireAuth, controllers.CreateCustomerContact)
	r.PUT("/contacts", middlewares.RequireAuth, controllers.UpdateCustomerContact)
	r.GET("/contacts/:uid", middlewares.RequireAuth, controllers.GetCustomerContact)
	////==== End contacts routes

	////==== Begin guarantors routes
	r.POST("/guarantors", middlewares.RequireAuth, controllers.CreateCustomerGuarantor)
	r.PUT("/guarantors", middlewares.RequireAuth, controllers.UpdateCustomerGuarantor)
	r.GET("/guarantors/:uid", middlewares.RequireAuth, controllers.GetCustomerGuarantor)
	////==== End guarantors routes

	////==== Begin referees routes
	r.POST("/referees", middlewares.RequireAuth, controllers.CreateCustomerReferee)
	r.PUT("/referees", middlewares.RequireAuth, controllers.UpdateCustomerReferee)
	r.GET("/referees/:uid", middlewares.RequireAuth, controllers.GetCustomerReferee)
	////==== End referees routes

	////==== Begin interactions routes
	r.GET("/interactions", middlewares.RequireAuth, controllers.GetCustomerConversations)
	////==== End interactions routes

	r.Run()
}
