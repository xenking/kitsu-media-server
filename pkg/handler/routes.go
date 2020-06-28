package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/xenking/kitsu-media-server/pkg/config"
	"github.com/xenking/kitsu-media-server/pkg/router/middleware"
)

func (h *Handler) Register(v1 *echo.Group) {
	v1.POST("/register", h.SignUp)
	v1.POST("/login", h.Login)

	jwtMiddleware := middleware.JWT(config.Global.JWTSecret)
	user := v1.Group("/user", jwtMiddleware)
	user.GET("", h.CurrentUser)
	user.PUT("", h.UpdateUser)

	users := v1.Group("/users", jwtMiddleware)
	users.GET("/:username", h.GetProfile)
	users.POST("/:username/follow", h.Follow)
	users.DELETE("/:username/follow", h.Unfollow)

	admin := v1.Group("/admin", jwtMiddleware)
	admin.GET("/users", h.UsersList)
	admin.DELETE("/user/:username", h.DeleteUser)

	articles := v1.Group("/articles", middleware.JWTWithConfig(
		middleware.JWTConfig{
			Skipper: func(c echo.Context) bool {
				if c.Request().Method == "GET" && c.Path() != "/api/articles/feed" {
					return true
				}
				return false
			},
			SigningKey: config.Global.JWTSecret,
		},
	))
	articles.POST("", h.CreateArticle)
	articles.GET("/feed", h.ArticleFeed)
	articles.PUT("/:slug", h.UpdateArticle)
	articles.DELETE("/:slug", h.DeleteArticle)
	articles.POST("/:slug/comments", h.AddArticleComment)
	articles.DELETE("/:slug/comments/:id", h.DeleteArticleComment)
	articles.POST("/:slug/favorite", h.ArticleFavorite)
	articles.DELETE("/:slug/favorite", h.ArticleUnfavorite)
	articles.GET("", h.Articles)
	articles.GET("/:slug", h.GetArticle)
	articles.GET("/:slug/comments", h.GetArticleComments)
	articleTags := articles.Group("/tags")
	articleTags.GET("", h.ArticleTags)

	medias := v1.Group("/medias", middleware.JWTWithConfig(
		middleware.JWTConfig{
			Skipper: func(c echo.Context) bool {
				if c.Request().Method == "GET" && c.Path() != "/api/medias/feed" {
					return true
				}
				return false
			},
			SigningKey: config.Global.JWTSecret,
		},
	))
	medias.POST("", h.CreateMedia)
	medias.GET("/feed", h.MediaFeed)
	medias.PUT("/:slug", h.UpdateMedia)
	medias.DELETE("/:slug", h.DeleteMedia)
	medias.POST("/:slug/comments", h.AddMediaComment)
	medias.DELETE("/:slug/comments/:id", h.DeleteMediaComment)
	medias.POST("/:slug/favorite", h.MediaFavorite)
	medias.DELETE("/:slug/favorite", h.MediaUnfavorite)
	medias.GET("", h.Medias)
	medias.GET("/:slug", h.GetMedia)
	medias.GET("/:slug/comments", h.GetMediaComments)

	mediaTags := medias.Group("/tags")
	mediaTags.GET("", h.MediaTags)
}
