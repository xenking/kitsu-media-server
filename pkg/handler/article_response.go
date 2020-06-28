package handler

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/xenking/kitsu-media-server/pkg/model"
	"github.com/xenking/kitsu-media-server/pkg/user"
)

type articleResponse struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"author"`
}

type singleArticleResponse struct {
	Article *articleResponse `json:"article"`
}

type articleListResponse struct {
	Articles      []*articleResponse `json:"articles"`
	ArticlesCount int                `json:"articlesCount"`
}

func newArticleResponse(c echo.Context, a *model.Article) *singleArticleResponse {
	ar := new(articleResponse)
	ar.TagList = make([]string, 0)
	ar.Slug = a.Slug
	ar.Title = a.Title
	ar.Description = a.Description
	ar.Body = a.Body
	ar.CreatedAt = a.CreatedAt
	ar.UpdatedAt = a.UpdatedAt
	for _, t := range a.Tags {
		ar.TagList = append(ar.TagList, t.Tag)
	}
	for _, u := range a.Favorites {
		if u.ID == userIDFromToken(c) {
			ar.Favorited = true
		}
	}
	ar.FavoritesCount = len(a.Favorites)
	ar.Author.Username = a.Author.Username
	ar.Author.Image = a.Author.Image
	ar.Author.Bio = a.Author.Bio
	ar.Author.Following = a.Author.FollowedBy(userIDFromToken(c))
	return &singleArticleResponse{ar}
}

func newArticleListResponse(us user.Store, userID uint, articles []model.Article, count int) *articleListResponse {
	r := new(articleListResponse)
	r.Articles = make([]*articleResponse, 0)
	for _, a := range articles {
		ar := new(articleResponse)
		ar.TagList = make([]string, 0)
		ar.Slug = a.Slug
		ar.Title = a.Title
		ar.Description = a.Description
		ar.Body = a.Body
		ar.CreatedAt = a.CreatedAt
		ar.UpdatedAt = a.UpdatedAt
		for _, t := range a.Tags {
			ar.TagList = append(ar.TagList, t.Tag)
		}
		for _, u := range a.Favorites {
			if u.ID == userID {
				ar.Favorited = true
			}
		}
		ar.FavoritesCount = len(a.Favorites)
		ar.Author.Username = a.Author.Username
		ar.Author.Image = a.Author.Image
		ar.Author.Bio = a.Author.Bio
		ar.Author.Following, _ = us.IsFollower(a.AuthorID, userID)
		r.Articles = append(r.Articles, ar)
	}
	r.ArticlesCount = count
	return r
}
