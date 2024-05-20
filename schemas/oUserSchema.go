package schemas

type UserSchema struct {
	UID       int    `json:"uid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	JoinDate  string `json:"joinDate"`
	UserGroup string `json:"userGroup"`
	Status    string `json:"status"`
}

type UserTokenSchema struct {
	Token     string   `json:"token"`
	ExpiresIn int      `json:"expiresIn"`
	TokenType string   `json:"tokenType"`
	Scope     []string `json:"scope"`
}

type UserLoginSchema struct {
	EmailOrPhone string `json:"emailOrPhone" binding:"required"`
	Password     string `json:"password" binding:"required"`
}
