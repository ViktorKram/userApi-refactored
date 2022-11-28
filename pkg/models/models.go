package models

import (
	"sync"
	"time"
)

type (
	User struct {
		CreatedAt   time.Time `json:"created_at"`
		DisplayName string    `json:"display_name"`
		Email       string    `json:"email"`
	}

	UserStore struct {
		sync.Mutex
		UserId int             `json:"userId"`
		List   map[string]User `json:"list"`
	}
)
