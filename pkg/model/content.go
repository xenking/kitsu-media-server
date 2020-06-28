package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Content struct {
	gorm.Model
	Slug     string `gorm:"unique_index;not null"`
	Title    string `gorm:"not null"`
	Author   User
	AuthorID uint
	Comments []Comment
}

type Article struct {
	Content
	Description string
	Body        string
	Favorites   []User `gorm:"many2many:article_favorites;"`
	Tags        []Tag  `gorm:"many2many:article_tags;association_autocreate:false"`
}

type Media struct {
	Content
	Description string
	Studio      string
	Episodes    int
	Type        string
	Poster      *string
	AiringDate  time.Time
	Favorites   []User `gorm:"many2many:media_favorites;"`
	Tags        []Tag  `gorm:"many2many:media_tags;association_autocreate:false"`
}
