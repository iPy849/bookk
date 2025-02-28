package models

import (
	"time"
)

type IUserService[T any] interface {
	GetUser(id string) (*T, error)
	GetUserBatch(ids []string) ([]*T, error)
	CreateUser(user *T) (string, error)
	CreateUserBatch(user []*T) ([]string, error)
	UpdateUser(user *T) error
	SoftDeleteUser(id string) (string, error)
	DeleteUser(id string) error
	DeleteUserBatch(id []string) error
	SetBan(id string, banUntil time.Time) error
}

const (
	ROLE_USER = iota
	ROLE_ADMIN
	ROLE_ENTERPRISE
)

type User struct {
	Id          string    `json:"user_id"`
	Email       string    `json:"email"`
	Role        int       `json:"user_role"`
	CreatedAt   time.Time `json:"created_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	BannedUntil time.Time `json:"banned_until"`
	LastAction  time.Time `json:"last_action"`
}

// TODO: Definir los métodos de interfaz
