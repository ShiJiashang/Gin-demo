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

// Login godoc
// @Summary 用户登录
// @Description 根据用户名和密码登录，成功后返回 JWT token。
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /login [post]
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	cfg := config.Load()
	loginResponse, err := service.Login(req, cfg.JWTSecret, cfg.JWTExpire)
	if err != nil {
		handleLoginError(c, err)
		return
	}

	response.Success(c, loginResponse)
}

// GetAuthorizedUsers godoc
// @Summary 获取鉴权用户列表
// @Description 需要 Authorization Bearer token，支持分页、搜索和排序。
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认 1"
// @Param size query int false "每页数量，默认 10，最大 100"
// @Param keyword query string false "按用户名模糊搜索"
// @Param sort query string false "排序字段：id、age、created_at"
// @Param order query string false "排序方向：asc 或 desc"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/userlist [post]
func GetAuthorizedUsers(c *gin.Context) {
	//鉴权通过中间件实现，这里直接返回用户列表
	GetUsers(c)
}

// GetUsers godoc
// @Summary 获取用户列表
// @Description 支持分页、搜索和排序。
// @Tags Users
// @Produce json
// @Param page query int false "页码，默认 1"
// @Param size query int false "每页数量，默认 10，最大 100"
// @Param keyword query string false "按用户名模糊搜索"
// @Param sort query string false "排序字段：id、age、created_at"
// @Param order query string false "排序方向：asc 或 desc"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
func GetUsers(c *gin.Context) {
	var query dto.UserListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	users, err := service.ListUsers(query)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "server error")
		return
	}

	response.Success(c, users)
}

// GetUser godoc
// @Summary 根据 ID 获取用户
// @Description 根据路径参数 id 获取一个用户。
// @Tags Users
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id} [get]
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

// CreateUser godoc
// @Summary 创建用户
// @Description 创建一个新用户。
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "创建用户请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [post]
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

// UpdateUser godoc
// @Summary 更新用户
// @Description 根据 id 更新用户姓名和年龄。
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "用户 ID"
// @Param request body dto.UpdateUserRequest true "更新用户请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id} [put]
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

// DeleteUser godoc
// @Summary 删除用户
// @Description 根据 id 删除用户。
// @Tags Users
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id} [delete]
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

func handleLoginError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrInvalidPassword) {
		response.Fail(c, http.StatusUnauthorized, 40104, "invalid name or password")
		return
	}

	response.Fail(c, http.StatusInternalServerError, 50001, "server error")
}
