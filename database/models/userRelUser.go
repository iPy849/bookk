package models

type UserRelUser struct {
	UserId        string `json:"user_id"`
	RelatedUserId string `json:"related_user_id"`
}

type IUserRelUserService interface {
	GetRelatedUsers(id string) ([]string, error)
	GetRelatedUsersByRole(id string, role int) ([]string, error)
	RelateUsers(userId, relatedUserId string) error
	UpdateUserRelation(userId, relatedUserId, updateId string) error
	RemoveRelation(userId, relatedUserId string) error
}
