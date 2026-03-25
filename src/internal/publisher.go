package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CooklyDev/AuthService/internal/adapters"
	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	logger           domain.Logger
	postgresPool     *pgxpool.Pool
	redisClient      *redis.Client
	sessionKeyPrefix string
	pollInterval     time.Duration
	batchSize        uint64
}

type outboxMessage struct {
	ID        uuid.UUID
	EventType string
	Payload   outboxPayload
}

type outboxPayload struct {
	SessionID  uuid.UUID `json:"session_id"`
	UserID     uuid.UUID `json:"user_id"`
	SessionKey string    `json:"session_key"`
	TTLSeconds int64     `json:"ttl_seconds"`
}

func NewPublisher(
	ctx context.Context,
	logger domain.Logger,
	postgresConfig *PostgresConfig,
	redisConfig *RedisConfig,
	appConfig *AppConfig,
) (*Publisher, error) {
	postgresPool, err := adapters.NewPostgresPool(
		ctx,
		logger,
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.DBName,
		postgresConfig.SSLMode,
	)
	if err != nil {
		return nil, err
	}

	redisClient, err := adapters.NewRedisClient(
		ctx,
		logger,
		redisConfig.Host,
		redisConfig.Port,
		redisConfig.Password,
	)
	if err != nil {
		postgresPool.Close()
		return nil, err
	}

	return &Publisher{
		logger:           logger,
		postgresPool:     postgresPool,
		redisClient:      redisClient,
		sessionKeyPrefix: appConfig.SessionPrefix,
		pollInterval:     appConfig.PublisherInterval,
		batchSize:        appConfig.PublisherBatchSize,
	}, nil
}

func (p *Publisher) Close() {
	if p.redisClient != nil {
		p.logger.Info("redis client closing: dependency=redis")
		if err := p.redisClient.Close(); err != nil {
			p.logger.Warn(
				fmt.Sprintf(
					"redis client close failed: dependency=redis error=%s",
					err.Error(),
				),
			)
		}
		p.redisClient = nil
		p.logger.Info("redis client closed: dependency=redis")
	}

	if p.postgresPool == nil {
		return
	}

	p.logger.Info("postgres pool closing: dependency=postgres")
	p.postgresPool.Close()
	p.postgresPool = nil
	p.logger.Info("postgres pool closed: dependency=postgres")
}

func (p *Publisher) Run(ctx context.Context) error {
	p.logger.Info(
		fmt.Sprintf(
			"publisher started: operation=run interval=%s batch_size=%d",
			p.pollInterval,
			p.batchSize,
		),
	)

	if err := p.processUnpublished(ctx); err != nil {
		return err
	}

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("publisher stopped: operation=run")
			return nil
		case <-ticker.C:
			if err := p.processUnpublished(ctx); err != nil {
				return err
			}
		}
	}
}

func (p *Publisher) processUnpublished(ctx context.Context) error {
	const query = `
		SELECT id, event_type, payload
		FROM outbox
		WHERE published = FALSE
		ORDER BY created_at, id
		LIMIT $1
	`

	rows, err := p.postgresPool.Query(ctx, query, p.batchSize)
	if err != nil {
		p.logger.Error(
			fmt.Sprintf(
				"publisher failed: operation=get_unpublished error=%s",
				err.Error(),
			),
		)

		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			message      outboxMessage
			payloadBytes []byte
		)

		if err := rows.Scan(&message.ID, &message.EventType, &payloadBytes); err != nil {
			p.logger.Error(
				fmt.Sprintf(
					"publisher failed: operation=scan_unpublished error=%s",
					err.Error(),
				),
			)

			return err
		}

		if err := json.Unmarshal(payloadBytes, &message.Payload); err != nil {
			p.logger.Error(
				fmt.Sprintf(
					"publisher failed: operation=unmarshal_payload outbox_id=%s error=%s",
					message.ID,
					err.Error(),
				),
			)

			return err
		}

		if err := p.publishMessage(ctx, message); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		p.logger.Error(
			fmt.Sprintf(
				"publisher failed: operation=iterate_unpublished error=%s",
				err.Error(),
			),
		)

		return err
	}

	return nil
}

func (p *Publisher) publishMessage(ctx context.Context, message outboxMessage) error {
	switch message.EventType {
	case "session.created":
		if err := p.redisClient.Set(
			ctx,
			p.sessionKey(message.Payload.SessionID, message.Payload.SessionKey),
			message.Payload.UserID.String(),
			time.Duration(message.Payload.TTLSeconds)*time.Second,
		).Err(); err != nil {
			p.logger.Error(
				fmt.Sprintf(
					"publisher failed: operation=add_session_to_redis outbox_id=%s session_id=%s user_id=%s error=%s",
					message.ID,
					message.Payload.SessionID,
					message.Payload.UserID,
					err.Error(),
				),
			)

			return err
		}
	case "session.deleted":
		if err := p.redisClient.Del(ctx, p.sessionKey(message.Payload.SessionID, message.Payload.SessionKey)).Err(); err != nil {
			p.logger.Error(
				fmt.Sprintf(
					"publisher failed: operation=delete_session_from_redis outbox_id=%s session_id=%s error=%s",
					message.ID,
					message.Payload.SessionID,
					err.Error(),
				),
			)

			return err
		}
	default:
		err := domain.NewBusinessRuleError("unsupported outbox event type")
		p.logger.Warn(
			fmt.Sprintf(
				"publisher failed: operation=publish_message outbox_id=%s event_type=%s error=%s",
				message.ID,
				message.EventType,
				err.Error(),
			),
		)

		return err
	}

	const markPublishedQuery = `
		UPDATE outbox
		SET published = TRUE
		WHERE id = $1
	`

	if _, err := p.postgresPool.Exec(ctx, markPublishedQuery, message.ID); err != nil {
		p.logger.Error(
			fmt.Sprintf(
				"publisher failed: operation=mark_published outbox_id=%s error=%s",
				message.ID,
				err.Error(),
			),
		)

		return err
	}

	p.logger.Info(
		fmt.Sprintf(
			"publisher completed: operation=publish_message outbox_id=%s event_type=%s",
			message.ID,
			message.EventType,
		),
	)

	return nil
}

func (p *Publisher) sessionKey(sessionID uuid.UUID, payloadSessionKey string) string {
	if payloadSessionKey != "" {
		return payloadSessionKey
	}

	return p.sessionKeyPrefix + sessionID.String()
}
