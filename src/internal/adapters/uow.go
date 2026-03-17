package adapters

import (
	"context"
	"fmt"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/CooklyDev/AuthService/internal/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWorkPostgres struct {
	pool        *pgxpool.Pool
	tx          DBTX
	logger      domain.Logger
	userRepo    usecases.UserRepo
	sessionRepo usecases.SessionRepo
}

func NewUnitOfWorkPostgres(pool *pgxpool.Pool, logger domain.Logger) *UnitOfWorkPostgres {
	uow := &UnitOfWorkPostgres{
		pool:   pool,
		logger: logger,
	}

	uow.bind(pool)

	return uow
}

func (u *UnitOfWorkPostgres) Begin() error {
	u.logger.Debug("uow: begin transaction")

	tx, err := u.pool.Begin(context.Background())
	if err != nil {
		u.logger.Error(
			fmt.Sprintf(
				"uow begin failed: operation=begin_tx error=%s",
				err.Error(),
			),
		)

		return err
	}

	u.bind(tx)

	return nil
}

func (u *UnitOfWorkPostgres) Commit() error {
	tx, ok := u.tx.(interface {
		Commit(ctx context.Context) error
	})
	if !ok {
		u.logger.Debug("uow: commit skipped (not a transaction)")
		return nil
	}

	u.logger.Debug("uow: commit transaction")

	err := tx.Commit(context.Background())
	if err != nil {
		u.logger.Error(
			fmt.Sprintf(
				"uow commit failed: operation=commit_tx error=%s",
				err.Error(),
			),
		)

		return err
	}

	u.bind(u.pool)

	return nil
}

func (u *UnitOfWorkPostgres) Rollback() error {
	tx, ok := u.tx.(interface {
		Rollback(ctx context.Context) error
	})
	if !ok {
		u.logger.Debug("uow: rollback skipped (not a transaction)")
		return nil
	}

	u.logger.Debug("uow: rollback transaction")

	err := tx.Rollback(context.Background())
	if err != nil {
		u.logger.Error(
			fmt.Sprintf(
				"uow rollback failed: operation=rollback_tx error=%s",
				err.Error(),
			),
		)

		return err
	}

	u.bind(u.pool)

	return nil
}

func (u *UnitOfWorkPostgres) UserRepository() usecases.UserRepo {
	return u.userRepo
}

func (u *UnitOfWorkPostgres) SessionRepository() usecases.SessionRepo {
	return u.sessionRepo
}

func (u *UnitOfWorkPostgres) bind(db DBTX) {
	u.tx = db
	u.userRepo = NewUserRepository(db)
	u.sessionRepo = NewSessionRepository(db)
}
