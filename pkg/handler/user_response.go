package handler

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/xenking/kitsu-media-server/pkg/model"
	"github.com/xenking/kitsu-media-server/pkg/user"
	"github.com/xenking/kitsu-media-server/pkg/utils"
)

type userResponse struct {
	User struct {
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Bio      *string `json:"bio"`
		Image    *string `json:"image"`
		Token    string  `json:"token"`
	} `json:"user"`
}

func newUserResponse(u *model.User) *userResponse {
	r := new(userResponse)
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Bio = u.Bio
	r.User.Image = u.Image
	r.User.Token = utils.GenerateJWT(u.ID)
	return r
}

type userListResponse struct {
	Users      []*userResponse `json:"users"`
	UsersCount int             `json:"usersCount"`
}

func newUserListResponse(users []model.User, count int) *userListResponse {
	r := new(userListResponse)
	r.Users = make([]*userResponse, 0)
	for _, a := range users {
		ur := new(userResponse)
		ur.User.Username = a.Username
		ur.User.Image = a.Image
		ur.User.Email = a.Email
		r.Users = append(r.Users, ur)
	}
	r.UsersCount = count
	return r
}

type profileResponse struct {
	Profile struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"profile"`
}

func newProfileResponse(us user.Store, userID uint, u *model.User) *profileResponse {
	r := new(profileResponse)
	r.Profile.Username = u.Username
	r.Profile.Bio = u.Bio
	r.Profile.Image = u.Image
	r.Profile.Following, _ = us.IsFollower(u.ID, userID)
	return r
}

type commentResponse struct {
	ID        uint      `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"author"`
}

type singleCommentResponse struct {
	Comment *commentResponse `json:"comment"`
}

type commentListResponse struct {
	Comments []commentResponse `json:"comments"`
}

func newCommentResponse(c echo.Context, cm *model.Comment) *singleCommentResponse {
	comment := new(commentResponse)
	comment.ID = cm.ID
	comment.Body = cm.Body
	comment.CreatedAt = cm.CreatedAt
	comment.UpdatedAt = cm.UpdatedAt
	comment.Author.Username = cm.User.Username
	comment.Author.Image = cm.User.Image
	comment.Author.Bio = cm.User.Bio
	comment.Author.Following = cm.User.FollowedBy(userIDFromToken(c))
	return &singleCommentResponse{comment}
}

func newCommentListResponse(c echo.Context, comments []model.Comment) *commentListResponse {
	r := new(commentListResponse)
	cr := commentResponse{}
	r.Comments = make([]commentResponse, 0)
	for _, i := range comments {
		cr.ID = i.ID
		cr.Body = i.Body
		cr.CreatedAt = i.CreatedAt
		cr.UpdatedAt = i.UpdatedAt
		cr.Author.Username = i.User.Username
		cr.Author.Image = i.User.Image
		cr.Author.Bio = i.User.Bio
		cr.Author.Following = i.User.FollowedBy(userIDFromToken(c))

		r.Comments = append(r.Comments, cr)
	}
	return r
}

type tagListResponse struct {
	Tags []string `json:"tags"`
}

func newTagListResponse(tags []model.Tag) *tagListResponse {
	r := new(tagListResponse)
	for _, t := range tags {
		r.Tags = append(r.Tags, t.Tag)
	}
	return r
}
