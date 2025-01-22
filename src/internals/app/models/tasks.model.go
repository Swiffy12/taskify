package models

import "github.com/jackc/pgtype"

type Task struct {
	Id          uint64             `json:"id"`
	Title       string             `valid:"required" json:"title"`
	Description string             `valid:"required" json:"description"`
	PriorityId  uint64             `valid:"required" json:"priority_id"`
	StatusId    uint64             `json:"status_id"`
	CreatedAt   pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
	CreatorId   uint64             `db:"creator_id" json:"creator_id"`
	AssignedId  uint64             `valid:"required" db:"assigned_id" json:"assigned_id"`
}

type CreateTaskRequestDTO struct {
	Title       string `valid:"required" json:"title"`
	Description string `valid:"required" json:"description"`
	PriorityId  uint64 `valid:"required, int" json:"priority_id"`
	AssignedId  uint64 `valid:"required, int" db:"assigned_id" json:"assigned_id"`
}

type GetTasksRequestDTO struct {
	Id         uint64
	Title      string
	CreatorId  uint64
	AssignedId uint64
}

type UpdateTaskRequestDTO struct {
	Title       string
	Description string
	PriorityId  uint64 `valid:"int" json:"priority_id"`
	StatusId    uint64 `valid:"int" json:"status_id"`
	CreatorId   uint64 `valid:"int" db:"creator_id" json:"creator_id"`
	AssignedId  uint64 `valid:"int" db:"assigned_id" json:"assigned_id"`
}
