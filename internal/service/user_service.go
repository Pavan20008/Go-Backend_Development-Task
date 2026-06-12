package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/Pavan20008/user-age-api/db/sqlc"
	"github.com/Pavan20008/user-age-api/internal/models"
	"github.com/Pavan20008/user-age-api/internal/repository"
)

// ErrNotFound is re-exported so callers (handlers) depend on the service layer
// rather than the repository directly.
var ErrNotFound = repository.ErrNotFound

// ErrInvalidDate is returned when a date of birth cannot be parsed.
var ErrInvalidDate = errors.New("invalid date of birth")

// UserService contains the business logic for managing users, including
// dynamic age calculation.
type UserService struct {
	repo   *repository.UserRepository
	logger *zap.Logger
}

// NewUserService builds a UserService.
func NewUserService(repo *repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func parseDOB(s string) (time.Time, error) {
	t, err := time.Parse(models.DateLayout, s)
	if err != nil {
		return time.Time{}, ErrInvalidDate
	}
	return t, nil
}

func (s *UserService) toResponse(u sqlc.User) models.UserResponse {
	return models.UserResponse{
		ID:   u.ID,
		Name: u.Name,
		Dob:  u.Dob.Time.Format(models.DateLayout),
	}
}

func (s *UserService) toResponseWithAge(u sqlc.User) models.UserWithAgeResponse {
	return models.UserWithAgeResponse{
		ID:   u.ID,
		Name: u.Name,
		Dob:  u.Dob.Time.Format(models.DateLayout),
		Age:  CalculateAge(u.Dob.Time, time.Now().UTC()),
	}
}

// Create stores a new user.
func (s *UserService) Create(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, err := parseDOB(req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("create user: %w", err)
	}

	s.logger.Info("user created", zap.Int32("id", user.ID), zap.String("name", user.Name))
	return s.toResponse(user), nil
}

// Get returns a user (with age) by id.
func (s *UserService) Get(ctx context.Context, id int32) (models.UserWithAgeResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.UserWithAgeResponse{}, err
	}
	return s.toResponseWithAge(user), nil
}

// List returns a page of users (with age).
func (s *UserService) List(ctx context.Context, limit, offset int32) ([]models.UserWithAgeResponse, error) {
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	out := make([]models.UserWithAgeResponse, 0, len(users))
	for _, u := range users {
		out = append(out, s.toResponseWithAge(u))
	}
	return out, nil
}

// Count returns the total number of users (used for pagination metadata).
func (s *UserService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Update modifies an existing user.
func (s *UserService) Update(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	dob, err := parseDOB(req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	s.logger.Info("user updated", zap.Int32("id", user.ID))
	return s.toResponse(user), nil
}

// Delete removes a user by id.
func (s *UserService) Delete(ctx context.Context, id int32) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.logger.Info("user deleted", zap.Int32("id", id))
	return nil
}
