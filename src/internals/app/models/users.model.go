package models

type User struct {
	Id       int64  `json:"id"`
	Fullname string `db:"full_name" json:"fullname`
	Rank     string `json:"rank"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
