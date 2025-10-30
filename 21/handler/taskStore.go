package handler

import (
	"restapi/task"
	"restapi/user"
)

type TaskStore interface {
	AddTask(t *task.Task) (*task.Task, error)
	GetTask(id int) (*task.Task, error)
	GetSelectedTasks(name, orderBy, sort string, limit *int) ([]task.Task, error)
	UpdateTask(t *task.Task) (*task.Task, error)
	DeleteTask(id int) error
	AddComment(taskID, author int, text string) (*task.Comment, error)
	InsertUser(data *user.UserData) (int, error)
	CheckUser(data *user.UserData) (int, error)
}
