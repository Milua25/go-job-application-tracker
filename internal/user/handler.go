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

func (h *UserHandler) RegisterRoutes(r gin.IRouter, authMiddleware, adminMiddleware gin.HandlerFunc) {
	g := r.Group("/users", authMiddleware)
	g.GET("/:id", h.GetUserByID)
	g.PATCH("/:id", h.UpdateUser)

	// Admin-only routes
	admin := g.Group("", adminMiddleware)
	admin.GET("", h.GetAllUsers)
	admin.POST("", h.CreateUser)
	admin.POST("/:id/*action", h.DeactivateReactivateUser)
	admin.DELETE("/:id", h.DeleteUser)
	admin.GET("/sessions", h.FindAllLoginUsers)
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

	if err := h.svc.deleteByID(c.Request.Context(), id); err != nil {
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

	u, err := h.svc.updateByID(c.Request.Context(), id, req)
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

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		render.ValidationError(c, err)
		return
	}

	// create the user
	u, err := h.svc.createUser(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailInUse):
			render.ConflictResponseError(c, "email already in use", err)
		default:
			slog.Error("failed to create user", "error", err)
			render.InternalServerError(c, "failed to create user", err)
		}
		return
	}

	slog.Info("user created successfully", "user_id", u.ID)
	render.OK(c, gin.H{
		"message": "user created successfully",
		"user":    toUserResp(u),
	})
}

func (h *UserHandler) DeactivateReactivateUser(c *gin.Context) {
	id := c.Param("id")
	action := c.Param("action") // Remove the leading slash from the action
	if len(action) < 2 {
		render.BadRequestError(c, "invalid action", errors.New("action must be 'deactivate' or 'reactivate'"))
		return
	}
	action = action[1:] // Remove the leading slash from the action
	isActive := true
	if _, err := uuid.Parse(id); err != nil {
		render.BadRequestError(c, "invalid user id", err)
		return
	}
	if action != "deactivate" && action != "reactivate" {
		slog.Info("invalid action for user", "user_id", id, "action", action)
		render.BadRequestError(c, "invalid action", errors.New("action must be 'deactivate' or 'reactivate'"))
		return
	}

	if action == "deactivate" {
		isActive = false
	}

	if _, err := h.svc.updateByID(c.Request.Context(), id, UpdateUserRequest{IsActive: &isActive}); err != nil {
		slog.Error("failed to update user", "user_id", id, "error", err)
		if errors.Is(err, ErrNotFound) {
			render.NotFoundError(c, "user not found")
			return
		}
		render.InternalServerError(c, "failed to update user", err)
		return
	}
	slog.Info("user "+action+"d"+" successfully", "user_id", id)
	render.OK(c, gin.H{
		"message": "user " + action + "d" + " successfully",
	})
}

func (h *UserHandler) FindAllLoginUsers(c *gin.Context) {
	users, err := h.svc.findAllWithSessions(c.Request.Context())
	if err != nil {
		slog.Error("failed to fetch users with sessions", "error", err)
		render.InternalServerError(c, "failed to fetch users with sessions", err)
		return
	}
	render.OK(c, gin.H{
		"users": toUserResps(users),
	})
}

// func setActiveStatus(c *gin.Context, isActive bool) {

// }
