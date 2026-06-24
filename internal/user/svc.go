package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

type userService struct {
	store Repository
}

func newUserService(store Repository) *userService {
	return &userService{store: store}
}

func (s *userService) getAll(ctx context.Context) ([]*User, error) {
	return s.store.GetAll(ctx)
}

func (s *userService) getByID(ctx context.Context, id string) (*User, error) {
	return s.store.GetByID(ctx, id)
}

func (s *userService) deleteByID(ctx context.Context, id string) error {
	if _, err := s.store.GetByID(ctx, id); err != nil {
		return err
	}
	return s.store.Delete(ctx, id)
}

func (s *userService) updateByID(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
	u, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		u.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		u.LastName = *req.LastName
	}
	if req.Email != nil {
		existing, err := s.store.GetByEmail(ctx, *req.Email)
		if err != nil && !errors.Is(err, ErrNotFound) {
			slog.Error("failed to check existing email", "email", *req.Email, "error", err)
			return nil, err
		}
		if existing != nil && existing.ID != u.ID {
			return nil, ErrEmailInUse
		}
		u.Email = *req.Email
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}

	if err := s.store.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// func (s *userService) getByEmail(ctx context.Context, email string) (*User, error) {
// 	return s.store.GetByEmail(ctx, email)
// }

func (s *userService) createUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	existingUser, err := s.store.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		slog.Error("failed to check existing user", "error", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailInUse
	}

	newUser := User{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		IsAdmin:   req.IsAdmin,
	}

	if err := s.store.Create(ctx, &newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *userService) findAllWithSessions(ctx context.Context) ([]*User, error) {
	return s.store.FindAllWithSessions(ctx)
}
