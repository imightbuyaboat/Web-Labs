package handler

import (
	"restapi/task"
)

type TaskCache interface {
	Get(taskID int) (*task.Task, error)
	Set(t *task.Task) error
	Delete(taskID int) error
}
