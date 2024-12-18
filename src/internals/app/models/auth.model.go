package models

type Auth struct {
	FullName string `valid:"required" json:"full_name"`
	Email    string `valid:"email, required" json:"email"`
	Password string `valid:"required" json:"password"`
}
