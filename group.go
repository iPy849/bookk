package bookk

import (
	"time"
)

type BaseGroup struct {
	Id        string
	Name      string
	CreatedAt time.Time
}

type Group struct {
	BaseGroup
	Description string
}

type IGroupService[G any, U any, I any] interface {
	GetGroupById(groupId string) (*G, error)
	CreateGroup(group G) (*G, error)
	UpdateGroup(group *G) error
	DeleteGroup(groupId string) error
	GetGroupUsers(groupId string) ([]*U, error)
	AddUserToGroup(groupId, userId string) error
	DeleteUserFromGroup(groupId, userId string) error
	GetGroupItems(groupId string) ([]*I, error)
	ExcludeGroupItem(groupId, itemId string) error
}
