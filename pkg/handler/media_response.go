package handler

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/xenking/kitsu-media-server/pkg/model"
	"github.com/xenking/kitsu-media-server/pkg/user"
)

type mediaResponse struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Studio         string    `json:"studio"`
	Episodes       int       `json:"episodes"`
	Type           string    `json:"type"`
	AiringDate     time.Time `json:"airingDate"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Poster         *string   `json:"poster"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"author"`
}

type singleMediaResponse struct {
	Media *mediaResponse `json:"media"`
}

type mediaListResponse struct {
	Medias      []*mediaResponse `json:"medias"`
	MediasCount int              `json:"mediasCount"`
}

func newMediaResponse(c echo.Context, m *model.Media) *singleMediaResponse {
	mr := new(mediaResponse)
	mr.TagList = make([]string, 0)
	mr.Slug = m.Slug
	mr.Title = m.Title
	mr.Description = m.Description
	mr.Studio = m.Studio
	mr.Episodes = m.Episodes
	mr.Type = m.Type
	mr.AiringDate = m.AiringDate
	mr.Poster = m.Poster
	mr.CreatedAt = m.CreatedAt
	mr.UpdatedAt = m.UpdatedAt
	for _, t := range m.Tags {
		mr.TagList = append(mr.TagList, t.Tag)
	}
	for _, u := range m.Favorites {
		if u.ID == userIDFromToken(c) {
			mr.Favorited = true
		}
	}
	mr.FavoritesCount = len(m.Favorites)
	mr.Author.Username = m.Author.Username
	mr.Author.Image = m.Author.Image
	mr.Author.Bio = m.Author.Bio
	mr.Author.Following = m.Author.FollowedBy(userIDFromToken(c))
	return &singleMediaResponse{mr}
}

func newMediaListResponse(us user.Store, userID uint, medias []model.Media, count int) *mediaListResponse {
	r := new(mediaListResponse)
	r.Medias = make([]*mediaResponse, 0)
	for _, m := range medias {
		mr := new(mediaResponse)
		mr.TagList = make([]string, 0)
		mr.Slug = m.Slug
		mr.Title = m.Title
		mr.Description = m.Description
		mr.Studio = m.Studio
		mr.Episodes = m.Episodes
		mr.Type = m.Type
		mr.AiringDate = m.AiringDate
		mr.Poster = m.Poster
		mr.CreatedAt = m.CreatedAt
		mr.UpdatedAt = m.UpdatedAt
		for _, t := range m.Tags {
			mr.TagList = append(mr.TagList, t.Tag)
		}
		for _, u := range m.Favorites {
			if u.ID == userID {
				mr.Favorited = true
			}
		}
		mr.FavoritesCount = len(m.Favorites)
		mr.Author.Username = m.Author.Username
		mr.Author.Image = m.Author.Image
		mr.Author.Bio = m.Author.Bio
		mr.Author.Following, _ = us.IsFollower(m.AuthorID, userID)
		r.Medias = append(r.Medias, mr)
	}
	r.MediasCount = count
	return r
}
