package adapters

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
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

type RedisSessionRepository struct {
	redisClient      *redis.Client
	sessionTTL       time.Duration
	sessionKeyPrefix string
}

func NewRedisSessionRepository(redisClient *redis.Client, sessionTTL time.Duration, sessionKeyPrefix string) *RedisSessionRepository {
	return &RedisSessionRepository{redisClient: redisClient, sessionTTL: sessionTTL, sessionKeyPrefix: sessionKeyPrefix}
}

func (r *RedisSessionRepository) Add(session *domain.Session) error {
	key := r.redisSessionKey(session.ID)
	value := session.UserID.String()

	err := r.redisClient.Set(context.Background(), key, value, r.sessionTTL).Err()
	if err != nil {
		return NewAdapterError("add session", err)
	}

	return nil
}

func (r *RedisSessionRepository) Delete(sessionID uuid.UUID) error {
	key := r.redisSessionKey(sessionID)

	err := r.redisClient.Del(context.Background(), key).Err()
	if err != nil {
		return NewAdapterError("delete session", err)
	}

	return nil
}

func (r *RedisSessionRepository) GetUserSessions(userID uuid.UUID) ([]*domain.Session, error) {
	ctx := context.Background()
	cursor := uint64(0)
	sessions := make([]*domain.Session, 0)

	for {
		keys, nextCursor, err := r.redisClient.Scan(ctx, cursor, r.sessionKeyPrefix+"*", 100).Result()
		if err != nil {
			return nil, NewAdapterError("get user sessions", err)
		}

		if len(keys) > 0 {
			values, err := r.redisClient.MGet(ctx, keys...).Result()
			if err != nil {
				return nil, NewAdapterError("get user sessions", err)
			}

			for index, key := range keys {
				if values[index] == nil {
					continue
				}

				storedUserID, err := redisUserID(values[index])
				if err != nil {
					return nil, NewAdapterError("get user sessions", err)
				}
				if storedUserID != userID {
					continue
				}

				sessionID, err := r.redisSessionID(key)
				if err != nil {
					return nil, NewAdapterError("get user sessions", err)
				}

				session, err := domain.NewSession(sessionID, storedUserID)
				if err != nil {
					return nil, NewAdapterError("get user sessions", err)
				}

				sessions = append(sessions, session)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return sessions, nil
}

func (r *RedisSessionRepository) GetSession(sessionID uuid.UUID) (*domain.Session, error) {
	value, err := r.redisClient.Get(context.Background(), r.redisSessionKey(sessionID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, NewAdapterError("get session", err)
	}

	userID, err := uuid.Parse(value)
	if err != nil {
		return nil, NewAdapterError("get session", err)
	}

	session, err := domain.NewSession(sessionID, userID)
	if err != nil {
		return nil, NewAdapterError("get session", err)
	}

	return session, nil
}

func (r *RedisSessionRepository) redisSessionKey(sessionID uuid.UUID) string {
	return r.sessionKeyPrefix + sessionID.String()
}

func (r *RedisSessionRepository) redisSessionID(key string) (uuid.UUID, error) {
	sessionID, ok := strings.CutPrefix(key, r.sessionKeyPrefix)
	if !ok || sessionID == "" {
		return uuid.Nil, fmt.Errorf("invalid session key %q", key)
	}

	return uuid.Parse(sessionID)
}

func redisUserID(value interface{}) (uuid.UUID, error) {
	switch typedValue := value.(type) {
	case string:
		return uuid.Parse(typedValue)
	case []byte:
		return uuid.Parse(string(typedValue))
	default:
		return uuid.Nil, fmt.Errorf("invalid session value type %T", value)
	}
}
