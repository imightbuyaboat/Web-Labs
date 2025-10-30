package db

import (
	"database/sql"
	"fmt"
	"restapi/task"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func (ps *PostgresStore) AddTask(t *task.Task) (*task.Task, error) {
	var insertedTask task.Task
	query := "insert into tasks (name, description) values ($1, $2) returning id, name, description"
	err := ps.db.QueryRow(query, t.Name, t.Description).Scan(&insertedTask.ID, &insertedTask.Name, &insertedTask.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to insert task: %v", err)
	}
	return &insertedTask, nil
}

func (ps *PostgresStore) GetTask(taskID int) (*task.Task, error) {
	var t task.Task
	query := `
		select 
		    t.id, t.name, t.description,
		    coalesce(
		        json_agg(
		            json_build_object(
		                'id', c.id,
		                'author', c.author,
		                'text', c.text,
		                'created_at', c.created_at
		            )
		        ) filter (where c.id is not null), '[]'
		    ) as comments
		from tasks t
		left join comments c on c.task_id = t.id
		where t.id = $1
		group by t.id;
	`

	err := ps.db.QueryRow(query, taskID).Scan(&t.ID, &t.Name, &t.Description, &t.Comments)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to select task %d from DB: %v", taskID, err)
	}
	return &t, nil
}

func (ps *PostgresStore) GetSelectedTasks(name, orderBy, sort string, limit *int) ([]task.Task, error) {
	query := `
		SELECT 
			t.id, t.name, t.description,
			COALESCE(
				json_agg(
					json_build_object(
						'id', c.id,
						'author', c.author,
						'text', c.text,
						'created_at', c.created_at
					)
				) FILTER (WHERE c.id IS NOT NULL),
				'[]'
			) AS comments
		FROM tasks t
		LEFT JOIN comments c ON c.task_id = t.id
	`
	var args []interface{}

	if name != "" {
		args = append(args, name)
		query += " where name = $" + strconv.Itoa(len(args))
	}

	query += " GROUP BY t.id, t.name, t.description"

	if orderBy != "" {
		query += " order by " + orderBy
		if strings.ToLower(sort) == "desc" {
			query += " desc"
		} else {
			query += " asc"
		}
	}

	if limit != nil {
		args = append(args, *limit)
		query += " limit $" + strconv.Itoa(len(args))
	}

	rows, err := ps.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select tasks from DB: %v", err)
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Comments); err != nil {
			return nil, fmt.Errorf("failed to scan task %d: %v", len(tasks)+1, err)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (ps *PostgresStore) UpdateTask(t *task.Task) (*task.Task, error) {
	var exists bool
	query := "select EXISTS (select 1 from tasks where id = $1)"
	err := ps.db.QueryRow(query, t.ID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if task %d exists: %v", t.ID, err)
	}
	if !exists {
		return nil, ErrTaskNotFound
	}

	query = "update tasks set name = $1, description = $2 where id = $3 returning id, name, description"
	var updatedTask task.Task

	err = ps.db.QueryRow(query, t.Name, t.Description, t.ID).Scan(&updatedTask.ID, &updatedTask.Name, &updatedTask.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to update task %d: %v", t.ID, err)
	}
	return &updatedTask, nil
}

func (ps *PostgresStore) DeleteTask(taskID int) error {
	query := "delete from tasks where id = $1"

	res, err := ps.db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task %d from DB: %v", taskID, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (ps *PostgresStore) AddComment(taskID, author int, text string) (*task.Comment, error) {
	query := `insert into comments (task_id, author, text) 
              values ($1, $2, $3) returning id, task_id, author, text, created_at`

	var c task.Comment
	err := ps.db.QueryRow(query, taskID, author, text).
		Scan(&c.ID, &c.TaskID, &c.Author, &c.Text, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert comment: %w", err)
	}

	return &c, nil
}
