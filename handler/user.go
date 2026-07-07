package handler

import (
	"errors"
	"net/http"
	"strconv"

	"gin_gorm_demo/config"
	"gin_gorm_demo/dto"
	"gin_gorm_demo/response"
	"gin_gorm_demo/service"

	"github.com/gin-gonic/gin"
)

// 登陆：根据用户名和密码验证用户是否存在，密码是否正确，返回token
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	config := config.Load()
	loginResponse, err := service.Login(req, config.JWTSecret, config.JWTExpire)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	response.Success(c, loginResponse)
}

// 鉴权token通过则返回用户列表
func GetAuthorizedUsers(c *gin.Context) {
	//鉴权通过中间件实现，这里直接返回用户列表
	GetUsers(c)
}

func GetUsers(c *gin.Context) {
	users, err := service.ListUsers()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "server error")
		return
	}

	response.Success(c, users)
}

func GetUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	user, err := service.GetUserByID(id)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response.Success(c, user)
}

func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	user, err := service.CreateUser(req)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response.Success(c, user)
}

func UpdateUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	user, err := service.UpdateUser(id, req)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response.Success(c, user)
}

func DeleteUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := service.DeleteUser(id); err != nil {
		handleUserError(c, err)
		return
	}

	response.Success(c, gin.H{
		"deleted": true,
	})
}

func parseID(c *gin.Context) (uint, bool) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40004, "invalid id")
		return 0, false
	}

	return uint(id64), true
}

func handleUserError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrUserNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "user not found")
		return
	}

	if errors.Is(err, service.ErrNameExists) {
		response.Fail(c, http.StatusConflict, 40901, "name already exists")
		return
	}

	if errors.Is(err, service.ErrInvalidAge) {
		response.Fail(c, http.StatusBadRequest, 40002, "invalid age")
		return
	}

	response.Fail(c, http.StatusInternalServerError, 50001, "server error")
}
