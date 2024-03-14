package users

import "time"

type CreateUserProfileRequest struct {
	Username string `json:"username" validate:"required,min=4,max=12"`
	FullName string `json:"fullName" validate:"required,max=36"`
}

type UserProfileResponse struct {
	// bson:"_id,omitempty" to remove it when we creating an item
	ID string `json:"id"`
	// json:"-" as we don't want to expose auth details
	Username  string    `json:"username"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
