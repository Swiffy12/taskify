package models

type User struct {
	Id           uint64  `json:"id"`
	FullName     string  `db:"full_name" json:"fullname"`
	Rank         *string `json:"rank"`
	Phone        *string `json:"phone"`
	Email        string  `json:"email"`
	PasswordHash string  `db:"password_hash"`
}

type GetUserResponseDTO struct {
	Id       uint64  `json:"id"`
	FullName string  `json:"fullname"`
	Rank     *string `json:"rank"`
	Phone    *string `json:"phone"`
	Email    string  `json:"email"`
}

type GetUsersRequestDTO struct {
	FullName string `json:"fullname"`
	Rank     string `json:"rank"`
}
