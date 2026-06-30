package user

import "time"

// UserResponse
type UserResponse struct {
	ID              string     `json:"id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Email           string     `json:"email"`
	IsActive        bool       `json:"is_active"`
	IsAdmin         bool       `json:"is_admin"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	IsAdmin   bool   `json:"is_admin" binding:"required"`
}

type DeactivateUserRequest struct {
	// Deactivate indicates whether the user should be deactivated (true) or activated (false).
	UserID     string `json:"user_id" binding:"required"`
	Deactivate bool   `json:"deactivate" binding:"required"`
}

func toUserResp(u *User) *UserResponse {
	return &UserResponse{
		ID:              u.ID.String(),
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		IsAdmin:         u.IsAdmin,
		Email:           u.Email,
		IsActive:        u.IsActive,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		LastLoginAt:     u.LastLoginAt,
		EmailVerifiedAt: u.EmailVerifiedAt,
	}
}

func toUserResps(users []*User) []UserResponse {
	resps := make([]UserResponse, len(users))
	for i, u := range users {
		resps[i] = *toUserResp(u)
	}
	return resps
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

type UpdateUserRoleRequest struct {
	IsAdmin bool `json:"is_admin"`
}

// func fromTimePtr(t *time.Time) *time.Time {
// 	if t == nil || t.IsZero() {
// 		return nil
// 	}
// 	return t
// }
