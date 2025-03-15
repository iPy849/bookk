package bookk

type BaseItem struct {
	Id          string
	UserId      string
	Name        string
	Description string
}

type Item struct {
	BaseItem
	price float32
}

type IItemRepository[T any] interface {
	GetItem(id string) (*T, error)
	GetItemBatch(id []string) ([]*T, error)
	CreateItem(item *T) (*T, error)
	CreateItemBatch(item []*T) ([]*T, error)
	UpdateItem(item *T) (*T, error)
	DeleteItem(id string) error
	DeleteItemBatch(id []string) error
	GetItemsByUserId(userId string) ([]*T, error)
	GetRelatedUserItems() ([]*T, error)
}
