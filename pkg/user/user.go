package user

import (
	"github.com/xenking/kitsu-media-server/pkg/model"
)

type Store interface {
	GetByID(uint) (*model.User, error)
	GetByEmail(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	Create(*model.User) error
	Update(*model.User) error
	Delete(*model.User) error
	List(offset, limit int) ([]model.User, int, error)

	AddFollower(user *model.User, followerID uint) error
	RemoveFollower(user *model.User, followerID uint) error
	IsFollower(userID, followerID uint) (bool, error)
}
