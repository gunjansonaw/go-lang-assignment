package models

import "time"

// User represents a user in the database
type User struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	DOB       time.Time `json:"dob" db:"dob"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserRequest represents the request payload for creating/updating a user
type UserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	DOB  string `json:"dob" validate:"required,date"`
}

// UserResponse represents the response payload for a user
type UserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  *int   `json:"age,omitempty"`
}

// UserListResponse represents paginated user list response
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
}

