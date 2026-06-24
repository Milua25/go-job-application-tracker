package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Milua25/go-job-application-tracker/internal/user"
	"github.com/Milua25/go-job-application-tracker/pkg/utils"
	"github.com/google/uuid"
)

func bootStrapAdmin(store user.Repository, adminEmail, adminPassword string) error {
	_, err := store.GetByEmail(context.Background(), adminEmail)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			// Admin user does not exist, create it
			// Hash the password
			hashedPassword, err := utils.HashPassword(adminPassword)
			if err != nil {
				return err
			}
			adminUser := user.User{
				ID:           uuid.New(),
				Email:        adminEmail,
				PasswordHash: hashedPassword,
				FirstName:    "Admin",
				LastName:     "User",
				IsAdmin:      true,
				IsActive:     true,
			}
			if err := store.Create(context.Background(), &adminUser); err != nil {
				return err
			}
			slog.Info("Admin user created successfully", "email", adminEmail)
			return nil
		}
		return err
	}
	slog.Info("Admin user already exists", "email", adminEmail)
	return nil // Admin user already exists
}
