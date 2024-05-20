package migrations

import (
	"super-lender/inits"
)

func init() {
	inits.LoadEnv()
	inits.DBInit()
}

func MigrateDatabase() {
	// inits.DB.AutoMigrate(&models.OCustomer{})
	// inits.DB.AutoMigrate(&models.OCustomerConversations{})
}
