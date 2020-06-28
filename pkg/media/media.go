package media

import (
	"github.com/xenking/kitsu-media-server/pkg/model"
)

type Store interface {
	GetBySlug(string) (*model.Media, error)
	GetUserMediaBySlug(userID uint, slug string) (*model.Media, error)
	CreateMedia(*model.Media) error
	UpdateMedia(*model.Media, []string) error
	DeleteMedia(*model.Media) error
	List(offset, limit int) ([]model.Media, int, error)
	ListByTag(tag string, offset, limit int) ([]model.Media, int, error)
	ListByAuthor(username string, offset, limit int) ([]model.Media, int, error)
	ListByWhoFavorited(username string, offset, limit int) ([]model.Media, int, error)
	ListFeed(userID uint, offset, limit int) ([]model.Media, int, error)

	AddComment(*model.Media, *model.Comment) error
	GetCommentsBySlug(string) ([]model.Comment, error)
	GetCommentByID(uint) (*model.Comment, error)
	DeleteComment(*model.Comment) error

	AddFavorite(*model.Media, uint) error
	RemoveFavorite(*model.Media, uint) error
	ListTags() ([]model.Tag, error)
}
