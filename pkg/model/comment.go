package model

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	Content   Content
	ContentID uint
	User      User
	UserID    uint
	Body      string
}

type Tag struct {
	gorm.Model
	Tag      string    `gorm:"unique_index"`
	Articles []Article `gorm:"many2many:article_tags;"`
	Medias   []Media   `gorm:"many2many:media_tags;"`
}
