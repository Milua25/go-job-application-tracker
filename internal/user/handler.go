package user

import (
	"errors"
	"log/slog"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	svc *userService
}

func NewUserHandler(store Repository) *UserHandler {
	return &UserHandler{svc: newUserService(store)}
}

func (h *UserHandler) RegisterRoutes(r gin.IRouter, authMiddleware gin.HandlerFunc) {
	g := r.Group("/users", authMiddleware)
	g.GET("", h.GetAllUsers)
	g.GET("/:id", h.GetUserByID)
	g.PATCH("/:id", h.UpdateUser)
	g.DELETE("/:id", h.DeleteUser)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	slog.Debug("fetching all users")

	users, err := h.svc.getAll(c.Request.Context())
	if err != nil {
		slog.Error("failed to fetch all users", "error", err)
		render.InternalServerError(c, "failed to get all users", err)
		return
	}

	slog.Debug("fetched users successfully", "count", len(users))
	render.OK(c, gin.H{
		"message": "get all users",
		"users":   toUserResps(users),
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	slog.Debug("fetching user by id", "user_id", id)

	if _, err := uuid.Parse(id); err != nil {
		render.BadRequestError(c, "invalid user id", err)
		return
	}

	u, err := h.svc.getByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			render.NotFoundError(c, "user not found")
			return
		}
		slog.Error("failed to fetch user by id", "user_id", id, "error", err)
		render.InternalServerError(c, "failed to get user by id", err)
		return
	}

	render.OK(c, gin.H{
		"message": "get user by id",
		"user":    toUserResp(u),
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	slog.Debug("deleting user", "user_id", id)

	if _, err := uuid.Parse(id); err != nil {
		render.BadRequestError(c, "invalid user id", err)
		return
	}

	if err := h.svc.delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			render.NotFoundError(c, "user not found")
			return
		}
		slog.Error("failed to delete user", "user_id", id, "error", err)
		render.InternalServerError(c, "failed to delete user", err)
		return
	}

	slog.Info("user deleted successfully", "user_id", id)
	render.NoContent(c)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		render.BadRequestError(c, "invalid user id", err)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		render.ValidationError(c, err)
		return
	}

	u, err := h.svc.update(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			render.NotFoundError(c, "user not found")
		case errors.Is(err, ErrEmailInUse):
			render.ConflictResponseError(c, "email already in use", err)
		default:
			slog.Error("failed to update user", "user_id", id, "error", err)
			render.InternalServerError(c, "failed to update user", err)
		}
		return
	}

	slog.Info("user updated successfully", "user_id", id)
	render.OK(c, gin.H{
		"message": "user updated successfully",
		"user":    toUserResp(u),
	})
}
