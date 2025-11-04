package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"restapi/task"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cache *redis.Client
	ctx   context.Context
}

func NewRedisCache() (*RedisCache, error) {
	rc := &RedisCache{}
	rc.ctx = context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err := client.Ping(rc.ctx).Err(); err != nil {
		return nil, err
	}

	rc.cache = client
	return rc, nil
}

func (rc *RedisCache) Set(t *task.Task) error {
	id := strconv.Itoa(t.ID)

	exists, err := rc.cache.HExists(rc.ctx, id, "name").Result()
	if err != nil {
		return fmt.Errorf("failed to check if task %d exists in cache: %v", t.ID, err)
	}
	if exists {
		return fmt.Errorf("task %d already exists in cache", t.ID)
	}

	err = rc.cache.HSet(rc.ctx, id, map[string]interface{}{
		"name":        t.Name,
		"description": t.Description,
		"comments":    t.Comments,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to insert task %d into cache: %v", t.ID, err)
	}

	err = rc.cache.Expire(rc.ctx, id, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set expired time of task %d in cache: %v", t.ID, err)
	}

	return nil
}

func (rc *RedisCache) Get(taskID int) (*task.Task, error) {
	id := strconv.Itoa(taskID)

	data, err := rc.cache.HGetAll(rc.ctx, id).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get task %d from cache: %v", taskID, err)
	}
	if len(data) == 0 {
		return nil, ErrTaskNotFound
	}

	t := &task.Task{
		ID:          taskID,
		Name:        data["name"],
		Description: data["description"],
		Comments:    json.RawMessage(data["comments"]),
	}

	return t, nil
}

func (rc *RedisCache) Delete(taskID int) error {
	id := strconv.Itoa(taskID)

	removed, err := rc.cache.Del(rc.ctx, id).Result()
	if err != nil {
		return fmt.Errorf("failed to delete task %d from cache: %v", taskID, err)
	}

	if removed == 0 {
		return ErrTaskNotFound
	}

	return nil
}
