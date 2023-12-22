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
type LoginSuccess struct {
	ID       string `json:"id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Token    string `json:"token,omitempty"`
}
type RegisterRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// NewDigimonHandler ...
func NewUserHandler(e *gin.Engine, userUsecase domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: userUsecase,
	}
	e.POST("/user/login", handler.Login)
	e.POST("/user/register", handler.Register)
}

// PostToCreateDigimon ...
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
	user.Token = token

	success(c, 0, "login success", LoginSuccess{
		ID:       user.ID,
		UserName: user.Name,
		Token:    user.Token,
	})
}

// PostToFosterDigimon ...
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
	}
	_, err = u.UserUsecase.FindByName(c, user.Name)

	if err == nil {
		success(c, USER_HAS_BEEN_REGISTER, "user has been register", nil)
		return
	}

	user, err = u.UserUsecase.Store(c, user)
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
