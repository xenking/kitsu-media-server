package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/xenking/kitsu-media-server/pkg/model"
	"github.com/xenking/kitsu-media-server/pkg/utils"
)

// GetMedia godoc
// @Summary Get an media
// @Description Get an media. Auth not required
// @ID get-media
// @ArticleTags media
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the media to get"
// @Success 200 {object} singleMediaResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /medias/{slug} [get]
func (h *Handler) GetMedia(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.mediaStore.GetBySlug(slug)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	return c.JSON(http.StatusOK, newMediaResponse(c, a))
}

// Medias godoc
// @Summary Get recent medias globally
// @Description Get most recent medias globally. Use query parameters to filter results. Auth is optional
// @ID get-medias
// @ArticleTags media
// @Accept  json
// @Produce  json
// @Param tag query string false "Filter by tag"
// @Param author query string false "Filter by author (username)"
// @Param favorited query string false "Filter by favorites of a user (username)"
// @Param limit query integer false "Limit number of medias returned (default is 20)"
// @Param offset query integer false "Offset/skip number of medias (default is 0)"
// @Success 200 {object} mediaListResponse
// @Failure 500 {object} utils.Error
// @Router /medias [get]
func (h *Handler) Medias(c echo.Context) error {
	var (
		medias []model.Media
		count  int
	)

	tag := c.QueryParam("tag")
	author := c.QueryParam("author")
	favoritedBy := c.QueryParam("favorited")

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 20
	}

	if tag != "" {
		medias, count, err = h.mediaStore.ListByTag(tag, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else if author != "" {
		medias, count, err = h.mediaStore.ListByAuthor(author, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else if favoritedBy != "" {
		medias, count, err = h.mediaStore.ListByWhoFavorited(favoritedBy, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else {
		medias, count, err = h.mediaStore.List(offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}

	return c.JSON(http.StatusOK, newMediaListResponse(h.userStore, userIDFromToken(c), medias, count))
}

// ArticleFeed godoc
// @Summary Get recent medias from users you follow
// @Description Get most recent medias from users you follow. Use query parameters to limit. Auth is required
// @ID feed
// @ArticleTags media
// @Accept  json
// @Produce  json
// @Param limit query integer false "Limit number of medias returned (default is 20)"
// @Param offset query integer false "Offset/skip number of medias (default is 0)"
// @Success 200 {object} mediaListResponse
// @Failure 401 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /medias/feed [get]
func (h *Handler) MediaFeed(c echo.Context) error {
	var (
		medias []model.Media
		count  int
	)

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 20
	}

	medias, count, err = h.mediaStore.ListFeed(userIDFromToken(c), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, newMediaListResponse(h.userStore, userIDFromToken(c), medias, count))
}

// CreateMedia godoc
// @Summary Create an media
// @Description Create an media. Auth is require
// @ID create-media
// @ArticleTags media
// @Accept  json
// @Produce  json
// @Param media body mediaCreateRequest true "Media to create"
// @Success 201 {object} singleMediaResponse
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /medias [post]
func (h *Handler) CreateMedia(c echo.Context) error {
	var a model.Media

	req := &mediaCreateRequest{}
	if err := req.bind(c, &a); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	a.AuthorID = userIDFromToken(c)

	err := h.mediaStore.CreateMedia(&a)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, newMediaResponse(c, &a))
}

// UpdateMedia godoc
// @Summary Update an media
// @Description Update an media. Auth is required
// @ID update-media
// @ArticleTags media
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the media to update"
// @Param media body mediaUpdateRequest true "Media to update"
// @Success 200 {object} singleMediaResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /medias/{slug} [put]
func (h *Handler) UpdateMedia(c echo.Context) error {
	slug := c.Param("slug")

	a, err := h.mediaStore.GetUserMediaBySlug(userIDFromToken(c), slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	req := &mediaUpdateRequest{}
	req.populate(a)

	if err := req.bind(c, a); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err = h.mediaStore.UpdateMedia(a, req.Media.Tags); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newMediaResponse(c, a))
}

// DeleteMedia godoc
// @Summary Delete an media
// @Description Delete an media. Auth is required
// @ID delete-media
// @ArticleTags media
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the media to delete"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /medias/{slug} [delete]
func (h *Handler) DeleteMedia(c echo.Context) error {
	slug := c.Param("slug")

	a, err := h.mediaStore.GetUserMediaBySlug(userIDFromToken(c), slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	err = h.mediaStore.DeleteMedia(a)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

// AddArticleComment godoc
// @Summary Create a comment for an media
// @Description Create a comment for an media. Auth is required
// @ID add-comment
// @ArticleTags comment
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the media that you want to create a comment for"
// @Param comment body createCommentRequest true "Comment you want to create"
// @Success 201 {object} singleCommentResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /medias/{slug}/comments [post]
func (h *Handler) AddMediaComment(c echo.Context) error {
	slug := c.Param("slug")

	a, err := h.mediaStore.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	var cm model.Comment

	req := &createCommentRequest{}
	if err := req.bind(c, &cm); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err = h.mediaStore.AddComment(a, &cm); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, newCommentResponse(c, &cm))
}

// GetArticleComments godoc
// @Summary Get the comments for an media
// @Description Get the comments for an media. Auth is optional
// @ID get-comments
// @ArticleTags comment
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the media that you want to get comments for"
// @Success 200 {object} commentListResponse
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /medias/{slug}/comments [get]
func (h *Handler) GetMediaComments(c echo.Context) error {
	slug := c.Param("slug")

	cm, err := h.mediaStore.GetCommentsBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newCommentListResponse(c, cm))
}

// DeleteArticleComment godoc
// @Summary Delete a comment for an media
// @Description Delete a comment for an media. Auth is required
// @ID delete-comments
// @ArticleTags comment
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the media that you want to delete a comment for"
// @Param id path integer true "ID of the comment you want to delete"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /medias/{slug}/comments/{id} [delete]
func (h *Handler) DeleteMediaComment(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	cm, err := h.mediaStore.GetCommentByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if cm == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	if cm.UserID != userIDFromToken(c) {
		return c.JSON(http.StatusUnauthorized, utils.NewError(errors.New("unauthorized action")))
	}

	if err := h.mediaStore.DeleteComment(cm); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

// ArticleFavorite godoc
// @Summary ArticleFavorite an media
// @Description ArticleFavorite an media. Auth is required
// @ID favorite
// @ArticleTags favorite
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the media that you want to favorite"
// @Success 200 {object} singleMediaResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /medias/{slug}/favorite [post]
func (h *Handler) MediaFavorite(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.mediaStore.GetBySlug(slug)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	if err := h.mediaStore.AddFavorite(a, userIDFromToken(c)); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newMediaResponse(c, a))
}

// ArticleUnfavorite godoc
// @Summary ArticleUnfavorite an media
// @Description ArticleUnfavorite an media. Auth is required
// @ID unfavorite
// @ArticleTags favorite
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the media that you want to unfavorite"
// @Success 200 {object} singleMediaResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /medias/{slug}/favorite [delete]
func (h *Handler) MediaUnfavorite(c echo.Context) error {
	slug := c.Param("slug")

	a, err := h.mediaStore.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	if err := h.mediaStore.RemoveFavorite(a, userIDFromToken(c)); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newMediaResponse(c, a))
}

// ArticleTags godoc
// @Summary Get tags
// @Description Get tags
// @ID tags
// @ArticleTags tag
// @Accept  json
// @Produce  json
// @Success 201 {object} tagListResponse
// @Failure 400 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tags [get]
func (h *Handler) MediaTags(c echo.Context) error {
	tags, err := h.mediaStore.ListTags()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, newTagListResponse(tags))
}
