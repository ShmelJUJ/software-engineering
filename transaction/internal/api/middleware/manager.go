package middleware

import (
	"fmt"
	"net/http"

	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ShmelJUJ/software-engineering/pkg/redis"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/generated/restapi/operations"
	"github.com/go-openapi/loads"
	"github.com/justinas/alice"
	"github.com/pwnedgod/idemgotent"
	"github.com/pwnedgod/wracha/adapter/goredis"
	"github.com/pwnedgod/wracha/logger/std"
)

// MiddlewareManager represents a manager for HTTP middlewares.
type MiddlewareManager struct {
	cfg   *Config
	log   logger.Logger
	redis *redis.Redis
	chain alice.Chain
}

// NewMiddlewareManager creates a new instance of MiddlewareManager.
func NewMiddlewareManager(cfg *Config, log logger.Logger, redis *redis.Redis) (*MiddlewareManager, error) {
	cfg, err := mergeWithDefault(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to set default config: %w", err)
	}

	return &MiddlewareManager{
		cfg:   cfg,
		log:   log,
		redis: redis,
		chain: alice.New(),
	}, nil
}

// AddIdempotenceMiddleware adds idempotence middleware to the middleware chain.
func (mm *MiddlewareManager) AddIdempotenceMiddleware() {
	middleware := idemgotent.Middleware(
		mm.cfg.IdempotencyCfg.Name,
		idemgotent.WithAdapter(goredis.NewAdapter(mm.redis.Client)),
		idemgotent.WithLogger(std.NewLogger()),
		idemgotent.WithKeySource(idemgotent.HeaderKeySource(mm.cfg.IdempotencyCfg.HeaderKey)),
	)

	mm.chain = mm.chain.Append(middleware)
}

// SetupGlobalMiddleware sets up global middleware based on Swagger specification.
func (mm *MiddlewareManager) SetupGlobalMiddleware(swaggerSpec *loads.Document, api *operations.TransactionAPI) {
	globalMiddleware := func(h http.Handler) http.Handler {
		return mm.chain.Then(h)
	}

	for path, pathInfo := range swaggerSpec.Spec().Paths.Paths {
		if pathInfo.Options != nil {
			api.AddMiddlewareFor("OPTIONS", path, globalMiddleware)
		}

		if pathInfo.Get != nil {
			api.AddMiddlewareFor("GET", path, globalMiddleware)
		}

		if pathInfo.Put != nil {
			api.AddMiddlewareFor("PUT", path, globalMiddleware)
		}

		if pathInfo.Post != nil {
			api.AddMiddlewareFor("POST", path, globalMiddleware)
		}

		if pathInfo.Delete != nil {
			api.AddMiddlewareFor("DELETE", path, globalMiddleware)
		}
	}
}
