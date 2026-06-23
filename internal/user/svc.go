package user

import (
	"context"
	"errors"
	"log/slog"
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

func (s *userService) delete(ctx context.Context, id string) error {
	if _, err := s.store.GetByID(ctx, id); err != nil {
		return err
	}
	return s.store.Delete(ctx, id)
}

func (s *userService) update(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
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
