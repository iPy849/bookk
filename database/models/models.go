package models

type AppModelProvider[UserType any, ItemType any] struct {
	UserModel    IUserService[UserType]
	UserRelation IUserRelationService
	ItemService  IItemService[ItemType]
}
