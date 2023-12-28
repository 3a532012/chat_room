package http

import (
	"log"
	"net/http"
	"user/domain"
	"utils/jwt"
	"utils/password"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS                = 0
	FAIL                   = 1
	USER_HAS_BEEN_REGISTER = 2
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}
type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
type AddFriendRequest struct {
	TagID string `json:"tag-id"`
}
type UpdateTagIDResponse struct {
	ID       string `json:"id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	TagID    string `json:"tag-id,omitempty"`
}
type FindUserByTagIDResponse struct {
	ID       string   `json:"id,omitempty"`
	UserName string   `json:"user_name,omitempty"`
	Friends  []string `json:"friends,omitempty"`
	TagID    string   `json:"tag-id,omitempty"`
}
type RemoveFriendRequest struct {
	TagID string `json:"tag-id"`
}
type LoginSuccess struct {
	ID       string `json:"id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Token    string `json:"token,omitempty"`
}
type RegisterRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func NewUserHandler(e *gin.Engine, userUsecase domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: userUsecase,
	}
	e.POST("/user/friend", jwt.AuthenticateJWT(), handler.AddFriend)
	e.DELETE("/user/friend", jwt.AuthenticateJWT(), handler.RemoveFriend)
	e.PATCH("/user/:tag-id", jwt.AuthenticateJWT(), handler.UpdateTagID)
	e.GET("/user/user/:tag-id", jwt.AuthenticateJWT(), handler.GetUserByTagID)
	e.POST("/user/login", handler.Login)
	e.POST("/user/register", handler.Register)
}
func (u *UserHandler) UpdateTagID(c *gin.Context) {
	tagID := c.Param("tag-id")
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
	updatedUser, err := u.UserUsecase.UpdateTagID(c, user.UserID, tagID)
	if err != nil {
		success(c, FAIL, err.Error(), nil)
		return
	}
	success(c, 0, "update tagID success", &UpdateTagIDResponse{
		ID:       updatedUser.ID,
		UserName: updatedUser.Name,
		TagID:    updatedUser.TagID,
	})
}
func (u *UserHandler) GetUserByTagID(c *gin.Context) {
	tagID := c.Param("tag-id")
	user, err := u.UserUsecase.FindByTagID(c, tagID)
	if err != nil {
		success(c, FAIL, err.Error(), nil)
		return
	}
	success(c, 0, "", &FindUserByTagIDResponse{
		ID:       user.ID,
		UserName: user.Name,
		Friends:  user.Friends,
		TagID:    user.TagID,
	})
}
func (u *UserHandler) AddFriend(c *gin.Context) {
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
	user, ok := claim.(*jwt.UserClaims)
	if !ok {
		success(c, FAIL, "jwt format fail", nil)
		return
	}
	err := u.UserUsecase.AddFriend(c, user.UserID, body.TagID)
	if err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}
	success(c, 0, "success add friend", nil)
}
func (u *UserHandler) RemoveFriend(c *gin.Context) {
	var body RemoveFriendRequest
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
	user, ok := claim.(*jwt.UserClaims)
	if !ok {
		success(c, FAIL, "jwt format fail", nil)
		return
	}
	err := u.UserUsecase.RemoveFriend(c, user.UserID, body.TagID)
	if err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}

	success(c, 0, "success remove friend", nil)
}
func (u *UserHandler) Login(c *gin.Context) {
	var body LoginRequest
	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}

	user, err := u.UserUsecase.FindByName(c, body.UserName)
	if err != nil {
		log.Println(err)
		success(c, FAIL, err.Error(), nil)
		return
	}
	if err := password.VerifyPassword(body.Password, user.Password); err != nil {
		success(c, FAIL, "password wrong", nil)
		return
	}

	token, err := jwt.CreateJWT(user.ID, user.Name)
	if err != nil {
		success(c, FAIL, err.Error(), nil)
		return
	}

	success(c, 0, "login success", LoginSuccess{
		ID:       user.ID,
		UserName: user.Name,
		Token:    token,
	})
}

func (u *UserHandler) Register(c *gin.Context) {
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

	user := &domain.User{
		Name:     body.UserName,
		Password: encryptPassword,
		Friends:  make([]string, 0),
	}
	_, err = u.UserUsecase.FindByName(c, user.Name)

	if err == nil {
		success(c, USER_HAS_BEEN_REGISTER, "user has been register", nil)
		return
	}

	_, err = u.UserUsecase.Store(c, user)
	// no find duplicate user

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
