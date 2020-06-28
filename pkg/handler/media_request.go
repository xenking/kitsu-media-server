package handler

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/xenking/kitsu-media-server/pkg/model"
	"time"
)

type mediaCreateRequest struct {
	Media struct {
		Title       string    `json:"title" validate:"required"`
		Description string    `json:"description" validate:"required"`
		Studio      string    `json:"studio" validate:"required"`
		Episodes    int       `json:"episodes" validate:"required"`
		Type        string    `json:"type" validate:"required"`
		AiringDate  time.Time `json:"airingDate"`
		Poster      string    `json:"poster"`
		Tags        []string  `json:"tagList, omitempty"`
	} `json:"media"`
}

func (r *mediaCreateRequest) bind(c echo.Context, m *model.Media) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	m.Title = r.Media.Title
	m.Slug = slug.Make(r.Media.Title)
	m.Description = r.Media.Description
	m.Studio = r.Media.Studio
	m.Episodes = r.Media.Episodes
	m.Type = r.Media.Type
	m.AiringDate = r.Media.AiringDate
	m.Poster = &r.Media.Poster
	if r.Media.Tags != nil {
		for _, t := range r.Media.Tags {
			m.Tags = append(m.Tags, model.Tag{Tag: t})
		}
	}
	return nil
}

type mediaUpdateRequest struct {
	Media struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Studio      string    `json:"studio"`
		Episodes    int       `json:"episodes"`
		Type        string    `json:"type"`
		AiringDate  time.Time `json:"airingDate"`
		Poster      string    `json:"poster"`
		Tags        []string  `json:"tagList"`
	} `json:"media"`
}

func (r *mediaUpdateRequest) populate(m *model.Media) {
	r.Media.Title = m.Title
	r.Media.Description = m.Description
	r.Media.Studio = m.Studio
	r.Media.Episodes = m.Episodes
	r.Media.Type = m.Type
	r.Media.AiringDate = m.AiringDate
	if m.Poster != nil {
		r.Media.Poster = *m.Poster
	}
}

func (r *mediaUpdateRequest) bind(c echo.Context, m *model.Media) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	m.Title = r.Media.Title
	m.Slug = slug.Make(m.Title)
	m.Description = r.Media.Description
	m.Studio = r.Media.Studio
	m.Episodes = r.Media.Episodes
	m.Type = r.Media.Type
	m.Poster = &r.Media.Poster
	m.AiringDate = r.Media.AiringDate
	return nil
}
