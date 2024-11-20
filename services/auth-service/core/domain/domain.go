package domain

import "time"

type Auth struct {
	ID                    string
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	RefreshTokenIsRevoked bool
	User                  User
	Claims                Claims
}

type User struct {
	ID       uint
	Username string
	Password string
	Email    string
	Role     string
}

type Claims struct {
	ID        uint
	Username  string
	Email     string
	Role      string
	Duration  time.Duration
	SessionID string
	IssuedAt  time.Time
	ExpiresAt time.Time
}
