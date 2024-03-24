package kafka

import (
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

var defaultLogger = watermill.NewStdLogger(false, false)

// RouterOption defines the functional option pattern for configuring the router.
type RouterOption func(*message.RouterConfig) error

// WithCloseTimeout sets the close timeout for the router.
func WithCloseTimeout(closeTimeout time.Duration) func(*message.RouterConfig) error {
	return func(rc *message.RouterConfig) error {
		rc.CloseTimeout = closeTimeout
		return nil
	}
}

// NewBrokerRouter creates a new message router with the given options.
func NewBrokerRouter(opts ...RouterOption) (*message.Router, error) {
	routerConfig := message.RouterConfig{}

	for _, opt := range opts {
		if err := opt(&routerConfig); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	router, err := message.NewRouter(routerConfig, defaultLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new router: %w", err)
	}

	return router, nil
}
