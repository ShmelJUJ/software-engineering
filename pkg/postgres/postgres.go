package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

// Default configurations for Postgres.
const (
	defaultMaxPoolSize = 1
	defaultConnAttemts = 10
	defaultConnTimeout = time.Second
)

type (
	// Option defines the functional option pattern for configuring Postgres.
	Option func(*Postgres) error

	// Postgres represents a Postgres instance with pool and transaction manager.
	Postgres struct {
		// Use Pool for better performance if there are no transactions.
		Pool *pgxpool.Pool
		// Use TrManager to work with transactions.
		TrManager *manager.Manager

		maxPoolSize  int
		connAttempts int
		connTimeout  time.Duration
		logger       logger.Logger
	}
)

// WithMaxPoolSize sets the maximum pool size.
func WithMaxPoolSize(maxSize int) Option {
	return func(p *Postgres) error {
		p.maxPoolSize = maxSize
		return nil
	}
}

// WithConnAttempts sets the number of connection attempts.
func WithConnAttempts(connAttempts int) Option {
	return func(p *Postgres) error {
		p.connAttempts = connAttempts
		return nil
	}
}

// WithConnTimeout sets the connection timeout.
func WithConnTimeout(connTimeout time.Duration) Option {
	return func(p *Postgres) error {
		p.connTimeout = connTimeout
		return nil
	}
}

// WithLogger sets the logger.
func WithLogger(log logger.Logger) Option {
	return func(p *Postgres) error {
		p.logger = log
		return nil
	}
}

// New creates a new Postgres instance with given url and functional options.
func New(ctx context.Context, url string, opts ...Option) (*Postgres, error) {
	defaultLogger, err := logger.NewLogrusLogger("info")
	if err != nil {
		return nil, fmt.Errorf("failed to create new default logger: %w", err)
	}

	p := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttemts,
		connTimeout:  defaultConnTimeout,
		logger:       defaultLogger,
	}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = int32(p.maxPoolSize)

	for p.connAttempts > 0 {
		p.Pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			break
		}

		p.logger.Info("Postgres is trying to connect...", map[string]interface{}{
			"attempts_left": p.connAttempts,
		})

		time.Sleep(p.connTimeout)

		p.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres pool: %w", err)
	}

	p.TrManager, err = manager.New(trmpgx.NewDefaultFactory(p.Pool))
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction manager: %w", err)
	}

	return p, nil
}

func (p *Postgres) GetTransactionConn(ctx context.Context) trmpgx.Tr {
	return trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, p.Pool)
}

// Close closes the Postgres pool.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
