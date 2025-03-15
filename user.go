package bookk

import (
	"time"
)

const (
	ROLE_USER = iota
	ROLE_ADVANCE_USER
	ROLE_ENTERPRISE
)

type BaseUser struct {
	Id          string
	Email       string
	Role        int
	CreatedAt   time.Time
	DeletedAt   time.Time
	BannedUntil time.Time
	LastAction  time.Time
}

type User struct {
	BaseUser
}

type IUserRepository[T any] interface {
	GetUser(id string) (*T, error)
	GetUserBatch(ids []string) ([]*T, error)
	CreateUser(user *T) (string, error)
	CreateUserBatch(user []*T) ([]string, error)
	UpdateUser(user *T) error
	DeleteUser(id string) error
	DeleteUserBatch(id []string) error
	SetBan(id string, banUntil time.Time) error
	GetRelatedUsers(id string) ([]*T, error)
	GetRelatedUsersByRole(id string, role int) ([]string, error)
	RelateUsers(userId, relatedUserId string) error
	RemoveRelation(userId, relatedUserId string) error
}
