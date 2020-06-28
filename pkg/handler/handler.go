package handler

import (
	"github.com/xenking/kitsu-media-server/pkg/article"
	"github.com/xenking/kitsu-media-server/pkg/media"
	"github.com/xenking/kitsu-media-server/pkg/user"
)

type Handler struct {
	userStore    user.Store
	articleStore article.Store
	mediaStore   media.Store
}

func NewHandler(us user.Store, as article.Store, ms media.Store) *Handler {
	return &Handler{
		userStore:    us,
		articleStore: as,
		mediaStore:   ms,
	}
}
