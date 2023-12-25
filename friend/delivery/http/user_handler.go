package http

import (
	"friend/domain"
	"log"
	"net/http"
	"utils/jwt"
	"utils/password"

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
type LoginRequest struct {
	FriendName string `json:"friend_name"`
	Password   string `json:"password"`
}
type LoginSuccess struct {
	ID         string `json:"id,omitempty"`
	FriendName string `json:"friend_name,omitempty"`
	Token      string `json:"token,omitempty"`
}
type RegisterRequest struct {
	FriendName string `json:"friend_name"`
	Password   string `json:"password"`
}

// NewDigimonHandler ...
func NewFriendHandler(e *gin.Engine, friendUsecase domain.FriendUsecase) {
	handler := &FriendHandler{
		FriendUsecase: friendUsecase,
	}
	e.POST("/friend/login", handler.Login)
	e.POST("/friend/register", handler.Register)
}

// PostToCreateDigimon ...
func (u *FriendHandler) Login(c *gin.Context) {
	var body LoginRequest
	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}

	friend, err := u.FriendUsecase.FindByName(c, body.FriendName)
	if err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}
	if err := password.VerifyPassword(body.Password, friend.Password); err != nil {
		success(c, FAIL, "password wrong", nil)
		return
	}

	token, err := jwt.CreateJWT(friend.ID, friend.Name)
	if err != nil {
		success(c, FAIL, err.Error(), nil)
		return
	}
	friend.Token = token

	success(c, 0, "login success", LoginSuccess{
		ID:         friend.ID,
		FriendName: friend.Name,
		Token:      friend.Token,
	})
}

// PostToFosterDigimon ...
func (u *FriendHandler) Register(c *gin.Context) {
	var body RegisterRequest
	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}
	encryptPassword, err := password.EncryptPassword(body.Password)

	if err != nil {
		success(c, USER_HAS_BEEN_REGISTER, err.Error(), nil)
		return
	}

	friend := &domain.Friend{
		Name:     body.FriendName,
		Password: encryptPassword,
	}
	_, err = u.FriendUsecase.FindByName(c, friend.Name)

	if err == nil {
		success(c, USER_HAS_BEEN_REGISTER, "friend has been register", nil)
		return
	}

	friend, err = u.FriendUsecase.Store(c, friend)
	// no find duplicate friend

	if err != nil {
		success(c, USER_HAS_BEEN_REGISTER, err.Error(), nil)
		return
	}

	success(c, SUCCESS, "Register OK", nil)
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
