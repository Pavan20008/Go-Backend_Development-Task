package models

// DateLayout is the canonical date format used across the API (ISO-8601 date).
const DateLayout = "2006-01-02"

// CreateUserRequest is the payload accepted by POST /users.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	Dob  string `json:"dob" validate:"required,dob"`
}

// UpdateUserRequest is the payload accepted by PUT /users/:id.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	Dob  string `json:"dob" validate:"required,dob"`
}

// UserResponse is returned by create and update endpoints (no age).
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
}

// UserWithAgeResponse is returned by get and list endpoints (includes age).
type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
	Age  int    `json:"age"`
}
