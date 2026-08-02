package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/murashi19/koda-b8-backend/internal/dto/users"
	"github.com/murashi19/koda-b8-backend/internal/lib/response"
	"github.com/murashi19/koda-b8-backend/internal/models"
	"github.com/murashi19/koda-b8-backend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GET /users
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success get users", users)
}

// GET /users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.GetDetailByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success get user", user)
}

// GET /users/me
func (h *UserHandler) Me(c *gin.Context) {

	userID := c.GetInt64("user_id")

	user, err := h.userService.GetMyProfile(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success get profile", user)
}

// PUT /users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var detail models.UserDetail

	if err := c.ShouldBindJSON(&detail); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	detail.ID = id

	if detail.Profile == nil {
		detail.Profile = &models.UserProfile{}
	}

	if err := h.userService.Update(c.Request.Context(), &detail); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User updated", nil)
}

// DELETE /users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}
	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "User deleted", nil)
}

// PATCH /users/me
func (h *UserHandler) UpdateMyProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req users.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.userService.UpdateMyProfile(
		c.Request.Context(),
		userID,
		&req,
	)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Profile updated successfully", nil)
}

// PATCH /users/me/avatar
func (h *UserHandler) UpdateAvatar(c *gin.Context) {
	userID := c.GetInt64("user_id")
	file, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "avatar is required")
		return
	}
	filename := strconv.FormatInt(userID, 10) + "_" + file.Filename
	path := "uploads/avatar/" + filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.userService.UpdateAvatar(
		c.Request.Context(),
		userID,
		path,
	); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Avatar updated successfully", gin.H{
		"avatar": path,
	})
}
