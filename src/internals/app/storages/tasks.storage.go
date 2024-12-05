package storages

import "github.com/jackc/pgx/v4/pgxpool"

type TasksStorage struct {
	databasePool *pgxpool.Pool
}

func NewTasksStorage(pool *pgxpool.Pool) *TasksStorage {
	tasksStorage := new(TasksStorage)
	tasksStorage.databasePool = pool
	return tasksStorage
}
