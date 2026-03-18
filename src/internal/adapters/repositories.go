package adapters

import (
	"context"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX represents either a connection or a transaction
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
}

// UserRepository implements application.UserRepo using PostgreSQL
type UserRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Add(user *domain.User) error {
	const query = `
		INSERT INTO users (id, username)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		user.ID,
		user.Username,
	)

	if err != nil {
		return NewAdapterError("add user", err)
	}

	return nil
}

// AuthIdentityRepository implements application.AuthIdentityRepo using PostgreSQL.
type AuthIdentityRepository struct {
	db DBTX
}

func NewAuthIdentityRepository(db DBTX) *AuthIdentityRepository {
	return &AuthIdentityRepository{db: db}
}

func (r *AuthIdentityRepository) Add(identity *domain.AuthIdentity) error {
	const query = `
		INSERT INTO auth_identities (id, user_id, provider, provider_id, email, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		identity.ID,
		identity.UserID,
		identity.Provider,
		identity.ProviderID,
		identity.Email,
		identity.PasswordHash,
	)
	if err != nil {
		return NewAdapterError("add auth identity", err)
	}

	return nil
}

func (r *AuthIdentityRepository) GetByEmail(email string) (*domain.AuthIdentity, error) {
	const query = `
		SELECT id, user_id, provider, provider_id, email, password_hash
		FROM auth_identities
		WHERE email = $1
	`

	row := r.db.QueryRow(context.Background(), query, email)

	var identity domain.AuthIdentity
	err := row.Scan(
		&identity.ID,
		&identity.UserID,
		&identity.Provider,
		&identity.ProviderID,
		&identity.Email,
		&identity.PasswordHash,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, NewAdapterError("get auth identity by email", err)
	}

	return &identity, nil
}

// SessionRepository implements application.SessionRepo using PostgreSQL
type SessionRepository struct {
	db DBTX
}

func NewSessionRepository(db DBTX) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Add(session *domain.Session) error {
	const query = `
		INSERT INTO sessions (id, user_id)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		session.ID,
		session.UserID,
	)

	if err != nil {
		return NewAdapterError("add session", err)
	}

	return nil
}

func (r *SessionRepository) Delete(sessionID string) error {
	const query = `
		DELETE FROM sessions
		WHERE id = $1
	`

	_, err := r.db.Exec(context.Background(), query, sessionID)
	if err != nil {
		return NewAdapterError("delete session", err)
	}

	return nil
}
