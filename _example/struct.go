package taskyapi

import "time"

// Error struct for error resource
type Error struct {
	Code        string `json:"code"`
	Detail      string `json:"detail"`
	ErrorFields []struct {
		Message string `json:"message,omitempty"`
		Name    string `json:"name,omitempty"`
	} `json:"errorFields,omitempty"`
}

// Task struct for task resource
type Task struct {
	CompletedAt time.Time `json:"completedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	ID          string    `json:"id"`
	Spent       int64     `json:"spent"`
	StartedAt   time.Time `json:"startedAt"`
	Status      string    `json:"status"`
	Tags        []string  `json:"tags"`
	Title       string    `json:"title"`
	User        *User     `json:"user,omitempty"`
}

// User struct for user resource
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TaskInstancesRequest struct for task
// GET: /tasks
type TaskInstancesRequest struct {
	Limit  int64 `json:"limit,omitempty" schema:"limit"`
	Offset int64 `json:"offset,omitempty" schema:"offset"`
}

// TaskInstancesResponse struct for task
// GET: /tasks
type TaskInstancesResponse []Task

// TaskSelfResponse struct for task
// GET: /tasks/{(#/definitions/task/definitions/identity)}
type TaskSelfResponse Task

// TaskCreateRequest struct for task
// POST: /tasks
type TaskCreateRequest struct {
	Tags  []string `json:"tags"`
	Title string   `json:"title"`
}

// TaskCreateResponse struct for task
// POST: /tasks
type TaskCreateResponse Task

// UserSelfResponse struct for user
// GET: /me
type UserSelfResponse User
