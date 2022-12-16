package model

type User struct {
	Username  string `json:"username" form:"username" example:"Username"`
	Password  string `json:"password" form:"password" example:"User password"`
	FirstName string `json:"first_name" form:"first_name" example:"First name"`
	LastName  string `json:"last_name" form:"last_name" example:"Last name"`
}
