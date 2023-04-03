package domain

import (
	"time"
)

const checkAuthTimer = 24 * 2 * time.Minute

type Auth struct {
	Domain     string    `yaml:"host"`
	Username   string    `yaml:"username"`
	UserID     int       `yaml:"userID"`
	AuthCookie string    `yaml:"authCookie"`
	Expires    time.Time `yaml:"expires"`
}
