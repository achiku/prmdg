package main

// User user resource
type User struct {
	ID        string `json:"id"`
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	Age       int    `json:"age"`
}

// Todo todo resource
type Todo struct {
	User  *User    `json:"user"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

// UserUpdateRequest user update response
type UserUpdateRequest *User

// UserUpdateResponse user update response
type UserUpdateResponse *User

// TodoCreateRequest todo create request
type TodoCreateRequest struct {
	UserID string   `json:"userId"`
	Title  string   `json:"title"`
	Tags   []string `json:"tags"`
}
