package handler

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"restapi/auth"
	"restapi/db"
	"restapi/middleware"
	"restapi/task"
	"restapi/user"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Handler struct {
	DB    TaskStore
	Cache TaskCache
}

func NewHandler(s TaskStore, c TaskCache) (*Handler, error) {
	return &Handler{
		DB:    s,
		Cache: c,
	}, nil
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var userData user.UserData

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID, err := h.DB.InsertUser(&userData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check user: %v", err), http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate JWT token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userData user.UserData

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID, err := h.DB.CheckUser(&userData)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, db.ErrIncorrectPassword) {
			http.Error(w, "Incorrect password", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to check user: %v", err), http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate JWT token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *Handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t task.Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if t.Name == "" || t.Description == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	insertedTask, err := h.DB.AddTask(&t)
	if err != nil {
		log.Printf("Failed to insert task into DB: %v", err)
		http.Error(w, fmt.Sprintf("Failed to insert task into DB: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedTask)
}

func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.Cache.Get(id)
	if err != nil {
		log.Printf("Failed to get from cache: %v", err)
	}

	if task == nil {
		task, err = h.DB.GetTask(id)
		if err != nil {
			if errors.Is(err, db.ErrTaskNotFound) {
				http.Error(w, fmt.Sprintf("Task %d not found", id), http.StatusNotFound)
			} else {
				log.Printf("Failed to get task from DB: %v", err)
				http.Error(w, fmt.Sprintf("Failed to get task from DB: %v", err), http.StatusInternalServerError)
			}
			return
		}

		if err = h.Cache.Set(task); err != nil {
			log.Printf("Failed to insert to cache: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

type GetSelectedTasksRequest struct {
	Name    string `json:"name,omitempty"`
	OrderBy string `json:"order_by,omitempty"`
	Sort    string `json:"sort,omitempty"`
	Limit   *int   `json:"limit,omitempty"`
	Format  string `json:"format,omitempty"`
}

func (h *Handler) GetSelectedTasksHandler(w http.ResponseWriter, r *http.Request) {
	var selectedTasksReq GetSelectedTasksRequest
	json.NewDecoder(r.Body).Decode(&selectedTasksReq)
	defer r.Body.Close()
	fmt.Println(selectedTasksReq)

	tasks, err := h.DB.GetSelectedTasks(
		selectedTasksReq.Name,
		selectedTasksReq.OrderBy,
		selectedTasksReq.Sort,
		selectedTasksReq.Limit,
	)
	if err != nil {
		log.Printf("Failed to get selected tasks from DB: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get selected tasks from DB: %v", err), http.StatusInternalServerError)
		return
	}

	if selectedTasksReq.Format == "" {
		selectedTasksReq.Format = "json"
	}

	switch strings.ToLower(selectedTasksReq.Format) {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=tasks.csv")
		w.WriteHeader(http.StatusOK)

		csvWriter := csv.NewWriter(w)
		defer csvWriter.Flush()

		if err := csvWriter.Write([]string{"ID", "Name", "Description"}); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write CSV header: %v", err), http.StatusInternalServerError)
			return
		}

		for _, t := range tasks {
			record := []string{
				strconv.Itoa(t.ID),
				t.Name,
				t.Description,
			}
			if err := csvWriter.Write(record); err != nil {
				http.Error(w, fmt.Sprintf("Failed to write CSV row: %v", err), http.StatusInternalServerError)
				return
			}
		}
	default:
		http.Error(w, "Unsupported format: "+selectedTasksReq.Format, http.StatusBadRequest)
	}
}

func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t task.Task

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	t.ID = id

	if t.ID == 0 || t.Name == "" || t.Description == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedTask, err := h.DB.UpdateTask(&t)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			http.Error(w, fmt.Sprintf("Task %d not found", t.ID), http.StatusNotFound)
		} else {
			log.Printf("Failed to update task in DB: %v", err)
			http.Error(w, fmt.Sprintf("Failed to update task in DB: %v", err), http.StatusInternalServerError)
		}
		return
	}

	if err = h.Cache.Delete(t.ID); err != nil {
		log.Printf("Failed to delete from cache: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.DB.DeleteTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			http.Error(w, fmt.Sprintf("Task %d not found", id), http.StatusNotFound)
		} else {
			log.Printf("Failed to delete task from DB: %v", err)
			http.Error(w, fmt.Sprintf("Failed to delete from DB: %v", err), http.StatusInternalServerError)
		}
		return
	}

	if err = h.Cache.Delete(id); err != nil {
		log.Printf("Failed to delete from cache: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AddCommentToTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var t struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	comment, err := h.DB.AddComment(id, userID, t.Text)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			http.Error(w, fmt.Sprintf("Task %d not found", id), http.StatusNotFound)
		} else {
			log.Printf("Failed to get task from DB: %v", err)
			http.Error(w, fmt.Sprintf("Failed to get task from DB: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}
