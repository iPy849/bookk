package models

import "time"

type Item struct {
	Id          string `json:"id"`
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type IItemService[T any] interface {
	GetItem(id string) (*T, error)
	GetItemBatch(id []string) ([]*T, error)
	CreateItem(item *T) (string, error)
	CreateItemBatch(item []*T) ([]string, error)
	UpdateItem(item *T) (string, error)
	DeleteItem(id string) (string, error)
	DeleteItemBatch(id []string) error
	GetItemsByUserId(userId string) ([]T, error)

	GetItemBookingsByDay(itemId string, day time.Time) ([]T, error)
	GetItemBookingsByTimeRange(itemId string, lowerBound time.Time, upperBound time.Time) ([]T, error)
}

type ItemRelUser struct{}
