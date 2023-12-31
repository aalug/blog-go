// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	UserID    int32     `json:"user_id"`
	PostID    int32     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	AuthorID    int32     `json:"author_id"`
	CategoryID  int32     `json:"category_id"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostTag struct {
	PostID int64 `json:"post_id"`
	TagID  int32 `json:"tag_id"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Tag struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	HashedPassword    string    `json:"hashed_password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	IsEmailVerified   bool      `json:"is_email_verified"`
}

type VerifyEmail struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}
