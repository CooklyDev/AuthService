package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

// DBTX represents either a connection or a transaction
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
}

// PgxUserRepository implements application.UserRepo using PostgreSQL.
type PgxUserRepository struct {
	db DBTX
}

func NewPgxUserRepository(db DBTX) *PgxUserRepository {
	return &PgxUserRepository{db: db}
}

func (r *PgxUserRepository) Add(user *domain.User) error {
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

// PgxAuthIdentityRepository implements application.AuthIdentityRepo using PostgreSQL.
type PgxAuthIdentityRepository struct {
	db DBTX
}

func NewPgxAuthIdentityRepository(db DBTX) *PgxAuthIdentityRepository {
	return &PgxAuthIdentityRepository{db: db}
}

func (r *PgxAuthIdentityRepository) Add(identity *domain.AuthIdentity) error {
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

func (r *PgxAuthIdentityRepository) GetByEmail(email string) (*domain.AuthIdentity, error) {
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

// PgxSessionRepository implements application.SessionRepo using PostgreSQL.
type PgxSessionRepository struct {
	db DBTX
}

func NewPgxSessionRepository(db DBTX) *PgxSessionRepository {
	return &PgxSessionRepository{db: db}
}

func (r *PgxSessionRepository) Add(session *domain.Session) error {
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

func (r *PgxSessionRepository) Delete(sessionID string) error {
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

type RedisSessionRepository struct {
	redisClient *redis.Client
	sessionTTL  time.Duration
}

func NewRedisSessionRepository(redisClient *redis.Client, sessionTTL time.Duration) *RedisSessionRepository {
	return &RedisSessionRepository{redisClient: redisClient, sessionTTL: sessionTTL}
}

func (r *RedisSessionRepository) Add(session *domain.Session) error {
	key := fmt.Sprintf("session:%s", session.ID.String())
	value := session.UserID.String()

	err := r.redisClient.Set(context.Background(), key, value, r.sessionTTL).Err()
	if err != nil {
		return NewAdapterError("add session", err)
	}

	return nil
}

func (r *RedisSessionRepository) Delete(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)

	err := r.redisClient.Del(context.Background(), key).Err()
	if err != nil {
		return NewAdapterError("delete session", err)
	}

	return nil
}
