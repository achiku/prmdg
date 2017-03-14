package exampleapi

import "time"

// Task struct for task resource
type Task struct {
	ID            string    `json:"id" schema:"id"`
	Identity      string    `json:"identity,omitempty" schema:"identity"`
	Name          string    `json:"name" schema:"name"`
	CreatedAt     time.Time `json:"createdAt" schema:"createdAt"`
	StartedAt     time.Time `json:"startedAt" schema:"startedAt"`
	UserID        string    `json:"userId,omitempty" schema:"userId"`
	Status        string    `json:"status" schema:"status"`
	TotalDuration int64     `json:"totalDuration" schema:"totalDuration"`
	CompletedAt   time.Time `json:"completedAt" schema:"completedAt"`
}

// User struct for user resource
type User struct {
	ID       string `json:"id" schema:"id"`
	Identity string `json:"identity,omitempty" schema:"identity"`
	Name     string `json:"name" schema:"name"`
}

// TaskSelfResponse struct for task
// GET: /tasks/{(#/definitions/task/definitions/identity)}
type TaskSelfResponse Task

// TaskCreateRequest struct for task
// POST: /tasks
type TaskCreateRequest struct {
	Name   string `json:"name" schema:"name"`
	UserID string `json:"userId" schema:"userId"`
}

// TaskCreateResponse struct for task
// POST: /tasks
type TaskCreateResponse Task

// TaskInstancesRequest struct for task
// GET: /tasks
type TaskInstancesRequest struct {
	Limit  int64 `json:"limit,omitempty" schema:"limit"`
	Offset int64 `json:"offset,omitempty" schema:"offset"`
}

// TaskInstancesResponse struct for task
// GET: /tasks
type TaskInstancesResponse []Task

// UserSelfResponse struct for user
// GET: /me
type UserSelfResponse User
