package http

import (
	"friend/domain"
	"log"
	"net/http"
	"utils/jwt"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS                = 0
	FAIL                   = 1
	USER_HAS_BEEN_REGISTER = 2
)

type FriendHandler struct {
	FriendUsecase domain.FriendUsecase
}
type AddFriendRequest struct {
	TagID string `json:"tag-id"`
}

func NewFriendHandler(e *gin.Engine, friendUsecase domain.FriendUsecase) {
	handler := &FriendHandler{
		FriendUsecase: friendUsecase,
	}
	e.POST("/friend", jwt.AuthenticateJWT(), handler.AddFriend)
	e.DELETE("/friend/:id", jwt.AuthenticateJWT(), handler.RemoveFriend)
	e.PATCH("friend/:id", jwt.AuthenticateJWT(), handler.AcceptFriend)
	e.GET("friend/list", jwt.AuthenticateJWT(), handler.FriendList)
	e.GET("/friend/request/list", jwt.AuthenticateJWT(), handler.FriendRequestList)
	e.GET("/friend/invite/list", jwt.AuthenticateJWT(), handler.FriendInviteList)
}
func (u *FriendHandler) AddFriend(c *gin.Context) {
	var body AddFriendRequest
	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}
	claim, isExist := c.Get("jwt")

	if !isExist {
		success(c, FAIL, "jwt doesn't in header", nil)
		return
	}
	sender, ok := claim.(*jwt.UserClaims)
	if !ok {
		success(c, FAIL, "jwt format fail", nil)
		return
	}
	err := u.FriendUsecase.AddFriend(c, sender.UserID, body.TagID)
	if err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}
	success(c, 0, "success add friend", nil)
}
func (u *FriendHandler) RemoveFriend(c *gin.Context) {
	id := c.Param("id")
	claim, isExist := c.Get("jwt")

	if !isExist {
		success(c, FAIL, "jwt doesn't in header", nil)
		return
	}
	user, ok := claim.(*jwt.UserClaims)
	if !ok {
		success(c, FAIL, "jwt format fail", nil)
		return
	}
	err := u.FriendUsecase.RemoveFriend(c, id, user.UserID)
	if err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}

	success(c, 0, "success remove friend", nil)
}

func (u *FriendHandler) AcceptFriend(c *gin.Context) {
	id := c.Param("id")
	claim, isExist := c.Get("jwt")

	if !isExist {
		success(c, FAIL, "jwt doesn't in header", nil)
		return
	}
	user, ok := claim.(*jwt.UserClaims)
	if !ok {
		success(c, FAIL, "jwt format fail", nil)
		return
	}
	err := u.FriendUsecase.AcceptFriend(c, id, user.UserID)
	if err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}

	success(c, 0, "success remove friend", nil)
}

func (u *FriendHandler) FriendList(c *gin.Context) {
	claim, isExist := c.Get("jwt")

	if !isExist {
		success(c, FAIL, "jwt doesn't in header", nil)
		return
	}
	user, ok := claim.(*jwt.UserClaims)
	if !ok {
		success(c, FAIL, "jwt format fail", nil)
		return
	}
	err := u.FriendUsecase.FriendList(c, user.UserID)
	if err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}

	success(c, 0, "success remove friend", nil)
}

func error(ctx *gin.Context, statusCode int, code int, message string) {
	ctx.JSON(statusCode, gin.H{
		"errorCode":    code,
		"errorMessage": message,
	})
}

func success(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}
