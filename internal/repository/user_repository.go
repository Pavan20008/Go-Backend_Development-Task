package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Pavan20008/user-age-api/db/sqlc"
)

// ErrNotFound is returned when a user does not exist.
var ErrNotFound = errors.New("user not found")

// UserRepository provides access to the users table via the sqlc-generated
// query layer. It translates between Go domain types and database types.
type UserRepository struct {
	q *sqlc.Queries
}

// NewUserRepository builds a UserRepository backed by the given querier source.
func NewUserRepository(db sqlc.DBTX) *UserRepository {
	return &UserRepository{q: sqlc.New(db)}
}

func toDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

// Create inserts a new user and returns the stored row.
func (r *UserRepository) Create(ctx context.Context, name string, dob time.Time) (sqlc.User, error) {
	return r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Name: name,
		Dob:  toDate(dob),
	})
}

// GetByID returns a single user by id, or ErrNotFound.
func (r *UserRepository) GetByID(ctx context.Context, id int32) (sqlc.User, error) {
	user, err := r.q.GetUser(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return sqlc.User{}, ErrNotFound
	}
	return user, err
}

// List returns a page of users ordered by id.
func (r *UserRepository) List(ctx context.Context, limit, offset int32) ([]sqlc.User, error) {
	return r.q.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

// Count returns the total number of users.
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}

// Update modifies an existing user and returns the updated row, or ErrNotFound.
func (r *UserRepository) Update(ctx context.Context, id int32, name string, dob time.Time) (sqlc.User, error) {
	user, err := r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  toDate(dob),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return sqlc.User{}, ErrNotFound
	}
	return user, err
}

// Delete removes a user by id, returning ErrNotFound when nothing was deleted.
func (r *UserRepository) Delete(ctx context.Context, id int32) error {
	rows, err := r.q.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
