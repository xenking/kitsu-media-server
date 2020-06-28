package store

import (
	"github.com/jinzhu/gorm"
	"github.com/xenking/kitsu-media-server/pkg/model"
)

type MediaStore struct {
	db *gorm.DB
}

func NewMediaStore(db *gorm.DB) *MediaStore {
	return &MediaStore{
		db: db,
	}
}

func (as *MediaStore) GetBySlug(s string) (*model.Media, error) {
	var m model.Media

	err := as.db.Where(&model.Media{Content: model.Content{Slug: s}}).Preload("Favorites").Preload("Tags").Preload("Author").Find(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &m, err
}

func (as *MediaStore) GetUserMediaBySlug(userID uint, slug string) (*model.Media, error) {
	var m model.Media

	err := as.db.Where(&model.Media{Content: model.Content{Slug: slug, AuthorID: userID}}).Find(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &m, err
}

func (as *MediaStore) CreateMedia(a *model.Media) error {
	tags := a.Tags

	tx := as.db.Begin()
	if err := tx.Create(&a).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, t := range a.Tags {
		err := tx.Where(&model.Tag{Tag: t.Tag}).First(&t).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&a).Association("Tags").Append(t).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Where(a.ID).Preload("Favorites").Preload("Tags").Preload("Author").Find(&a).Error; err != nil {
		tx.Rollback()
		return err
	}

	a.Tags = tags

	return tx.Commit().Error
}

func (as *MediaStore) UpdateMedia(a *model.Media, tagList []string) error {
	tx := as.db.Begin()
	if err := tx.Model(a).Update(a).Error; err != nil {
		tx.Rollback()
		return err
	}

	tags := make([]model.Tag, 0)

	for _, t := range tagList {
		tag := model.Tag{Tag: t}

		err := tx.Where(&tag).First(&tag).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			return err
		}

		tags = append(tags, tag)
	}

	if err := tx.Model(a).Association("Tags").Replace(tags).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where(a.ID).Preload("Favorites").Preload("Tags").Preload("Author").Find(a).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (as *MediaStore) DeleteMedia(a *model.Media) error {
	return as.db.Delete(a).Error
}

func (as *MediaStore) List(offset, limit int) ([]model.Media, int, error) {
	var (
		articles []model.Media
		count    int
	)

	as.db.Model(&articles).Count(&count)
	as.db.Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").Find(&articles)

	return articles, count, nil
}

func (as *MediaStore) ListByTag(tag string, offset, limit int) ([]model.Media, int, error) {
	var (
		t        model.Tag
		articles []model.Media
		count    int
	)

	err := as.db.Where(&model.Tag{Tag: tag}).First(&t).Error
	if err != nil {
		return nil, 0, err
	}

	as.db.Model(&t).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Association("Medias").
		Find(&articles)

	count = as.db.Model(&t).Association("Medias").Count()

	return articles, count, nil
}

func (as *MediaStore) ListByAuthor(username string, offset, limit int) ([]model.Media, int, error) {
	var (
		u        model.User
		articles []model.Media
		count    int
	)

	err := as.db.Where(&model.User{Username: username}).First(&u).Error
	if err != nil {
		return nil, 0, err
	}

	as.db.Where(&model.Media{Content: model.Content{AuthorID: u.ID}}).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&articles)
	as.db.Where(&model.Media{Content: model.Content{AuthorID: u.ID}}).Model(&model.Media{}).Count(&count)

	return articles, count, nil
}

func (as *MediaStore) ListByWhoFavorited(username string, offset, limit int) ([]model.Media, int, error) {
	var (
		u        model.User
		articles []model.Media
		count    int
	)

	err := as.db.Where(&model.User{Username: username}).First(&u).Error
	if err != nil {
		return nil, 0, err
	}

	as.db.Model(&u).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Association("MediaFavorites").
		Find(&articles)

	count = as.db.Model(&u).Association("MediaFavorites").Count()

	return articles, count, nil
}

func (as *MediaStore) ListFeed(userID uint, offset, limit int) ([]model.Media, int, error) {
	var (
		u        model.User
		articles []model.Media
		count    int
	)

	err := as.db.First(&u, userID).Error
	if err != nil {
		return nil, 0, err
	}

	var followings []model.Follow

	as.db.Model(&u).Preload("Following").Preload("Follower").Association("Followings").Find(&followings)

	if len(followings) == 0 {
		return articles, 0, nil
	}

	ids := make([]uint, len(followings))
	for i, f := range followings {
		ids[i] = f.FollowingID
	}

	as.db.Where("author_id in (?)", ids).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&articles)
	as.db.Where(&model.Media{Content: model.Content{AuthorID: u.ID}}).Model(&model.Media{}).Count(&count)

	return articles, count, nil
}

func (as *MediaStore) AddComment(a *model.Media, c *model.Comment) error {
	err := as.db.Model(a).Association("Comments").Append(c).Error
	if err != nil {
		return err
	}

	return as.db.Where(c.ID).Preload("User").First(c).Error
}

func (as *MediaStore) GetCommentsBySlug(slug string) ([]model.Comment, error) {
	var m model.Media
	err := as.db.Where(&model.Media{Content: model.Content{Slug: slug}}).Preload("Comments").Preload("Comments.User").First(&m).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return m.Comments, nil
}

func (as *MediaStore) GetCommentByID(id uint) (*model.Comment, error) {
	var m model.Comment
	if err := as.db.Where(id).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &m, nil
}

func (as *MediaStore) DeleteComment(c *model.Comment) error {
	return as.db.Delete(c).Error
}

func (as *MediaStore) AddFavorite(a *model.Media, userID uint) error {
	usr := model.User{}
	usr.ID = userID

	return as.db.Model(a).Association("Favorites").Append(&usr).Error
}

func (as *MediaStore) RemoveFavorite(a *model.Media, userID uint) error {
	usr := model.User{}
	usr.ID = userID

	return as.db.Model(a).Association("Favorites").Delete(&usr).Error
}

func (as *MediaStore) ListTags() ([]model.Tag, error) {
	var tags []model.Tag
	if err := as.db.Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}
