package models

import "github.com/jackc/pgx/pgtype"

type Task struct {
	Id          int64            `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Priority    uint8            `json:"priority"`
	Status      uint8            `json:"status"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"createdAt"`
	UpdatedAt   pgtype.Timestamp `db:"updated_at" json:"updatedAt"`
	CreatorId   uint8            `db:"creator_id" json:"creatorId"`
	AssignedId  uint8            `db:"assigned_id" json:"assignedId"`
}
