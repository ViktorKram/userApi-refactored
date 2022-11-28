package models

import (
	"time"
)

type (
	User struct {
		CreatedAt   time.Time `json:"created_at"`
		DisplayName string    `json:"display_name"`
		Email       string    `json:"email"`
	}

	UserStore struct {
		UserId int             `json:"userId"`
		List   map[string]User `json:"list"`
	}
)
