package models

import "github.com/jackc/pgx/pgtype"

type Task struct {
	Id          int64            `json:"id"`
	Title       string           `valid:"required" json:"title"`
	Description string           `valid:"required" json:"description"`
	Priority    uint8            `valid:"required" json:"priority"`
	Status      uint8            `json:"status"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt   pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	CreatorId   uint8            `db:"creator_id" json:"creator_id"`
	AssignedId  uint8            `valid:"required" db:"assigned_id" json:"assigned_id"`
}

type UpdateTaskRequestDTO struct {
	Id          int64 `valid:"int"`
	Title       string
	Description string
	Priority    uint8 `valid:"int"`
	Status      uint8 `valid:"int"`
	CreatorId   uint8 `valid:"int" db:"creator_id" json:"creator_id"`
	AssignedId  uint8 `valid:"int" db:"assigned_id" json:"assigned_id"`
}
