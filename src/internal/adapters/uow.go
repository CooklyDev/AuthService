package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/CooklyDev/AuthService/internal/application"
	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type UnitOfWorkApp struct {
	pool             *pgxpool.Pool
	redisClient      *redis.Client
	pgxTX            DBTX
	logger           domain.Logger
	userRepo         application.UserRepo
	authIdentityRepo application.AuthIdentityRepo
	sessionRepo      application.SessionRepo
	sessionTTL       time.Duration
	sessionKeyPrefix string
}

func NewUnitOfWorkApp(
	pool *pgxpool.Pool,
	client *redis.Client,
	logger domain.Logger,
	sessionTTL time.Duration,
	sessionKeyPrefix string,
) *UnitOfWorkApp {
	uow := &UnitOfWorkApp{
		pool:             pool,
		redisClient:      client,
		logger:           logger,
		sessionTTL:       sessionTTL,
		sessionKeyPrefix: sessionKeyPrefix,
	}

	uow.bind(pool, client)

	return uow
}

func (u *UnitOfWorkApp) Begin() error {
	u.logger.Debug("uow: begin transaction")

	pgxTX, err := u.pool.Begin(context.Background())
	if err != nil {
		u.logger.Error(
			fmt.Sprintf(
				"uow begin failed: operation=begin_tx error=%s",
				err.Error(),
			),
		)

		return err
	}
	u.bind(pgxTX, u.redisClient)

	return nil
}

func (u *UnitOfWorkApp) Commit() error {
	pgxTX, ok := u.pgxTX.(interface {
		Commit(ctx context.Context) error
	})
	if !ok {
		u.logger.Debug("uow: commit skipped (not a transaction)")
		return nil
	}

	u.logger.Debug("uow: commit transaction")

	err := pgxTX.Commit(context.Background())
	if err != nil {
		u.logger.Error(
			fmt.Sprintf(
				"uow commit failed: operation=commit_tx error=%s",
				err.Error(),
			),
		)

		return err
	}

	u.bind(u.pool, u.redisClient)

	return nil
}

func (u *UnitOfWorkApp) Rollback() error {
	pgxTX, ok := u.pgxTX.(interface {
		Rollback(ctx context.Context) error
	})
	if !ok {
		u.logger.Debug("uow: rollback skipped (not a transaction)")
		return nil
	}

	u.logger.Debug("uow: rollback transaction")

	err := pgxTX.Rollback(context.Background())
	if err != nil {
		u.logger.Error(
			fmt.Sprintf(
				"uow rollback failed: operation=rollback_tx error=%s",
				err.Error(),
			),
		)

		return err
	}

	u.bind(u.pool, u.redisClient)

	return nil
}

func (u *UnitOfWorkApp) UserRepository() application.UserRepo {
	return u.userRepo
}

func (u *UnitOfWorkApp) AuthIdentityRepository() application.AuthIdentityRepo {
	return u.authIdentityRepo
}

func (u *UnitOfWorkApp) SessionRepository() application.SessionRepo {
	return u.sessionRepo
}

func (u *UnitOfWorkApp) bind(db DBTX, redis *redis.Client) {
	u.pgxTX = db
	u.redisClient = redis
	u.userRepo = NewPgxUserRepository(db)
	u.authIdentityRepo = NewPgxAuthIdentityRepository(db)
	u.sessionRepo = NewRedisSessionRepository(redis, u.sessionTTL, u.sessionKeyPrefix)
}
