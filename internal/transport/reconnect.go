package transport

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// ReconnectConfig configures automatic reconnection behavior.
type ReconnectConfig struct {
	// Enabled turns auto-reconnect on or off.
	Enabled bool
	// MaxAttempts is the maximum number of reconnection attempts (0 = unlimited).
	MaxAttempts int
	// BaseWait is the initial wait duration between reconnection attempts.
	BaseWait time.Duration
	// MaxWait is the maximum wait duration between reconnection attempts.
	MaxWait time.Duration
	// OnReconnect is called after a successful reconnection.
	OnReconnect func() error
}

// Reconnector manages automatic reconnection for a WSTransport.
type Reconnector struct {
	transport *WSTransport
	config    ReconnectConfig

	mu         sync.Mutex
	active     bool
	cancelFunc context.CancelFunc
}

// NewReconnector creates a new reconnector for the given transport.
func NewReconnector(transport *WSTransport, config ReconnectConfig) *Reconnector {
	if config.BaseWait <= 0 {
		config.BaseWait = 1 * time.Second
	}
	if config.MaxWait <= 0 {
		config.MaxWait = 30 * time.Second
	}
	return &Reconnector{
		transport: transport,
		config:    config,
	}
}

// Start begins monitoring and reconnecting when disconnected.
// It should be called after the initial connection is established.
func (r *Reconnector) Start(ctx context.Context) {
	if !r.config.Enabled {
		return
	}

	r.mu.Lock()
	if r.active {
		r.mu.Unlock()
		return
	}
	r.active = true
	childCtx, cancel := context.WithCancel(ctx)
	r.cancelFunc = cancel
	r.mu.Unlock()

	go r.monitor(childCtx)
}

// Stop stops the reconnector.
func (r *Reconnector) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cancelFunc != nil {
		r.cancelFunc()
		r.cancelFunc = nil
	}
	r.active = false
}

// TriggerReconnect initiates a reconnection attempt.
func (r *Reconnector) TriggerReconnect(ctx context.Context) error {
	if !r.config.Enabled {
		return fmt.Errorf("reconnection is disabled")
	}
	return r.reconnect(ctx)
}

func (r *Reconnector) monitor(ctx context.Context) {
	// The monitor doesn't poll; instead, reconnection is triggered
	// by the OnDisconnect callback in the WSHandler.
	<-ctx.Done()
}

func (r *Reconnector) reconnect(ctx context.Context) error {
	for attempt := 1; r.config.MaxAttempts == 0 || attempt <= r.config.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Calculate exponential backoff with jitter.
		wait := time.Duration(float64(r.config.BaseWait) * math.Pow(2, float64(attempt-1)))
		if wait > r.config.MaxWait {
			wait = r.config.MaxWait
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}

		// Attempt to reconnect.
		if err := r.transport.Connect(ctx); err != nil {
			continue
		}

		// Run post-reconnect callback (e.g., re-login, re-subscribe).
		if r.config.OnReconnect != nil {
			if err := r.config.OnReconnect(); err != nil {
				// Close the connection and try again.
				_ = r.transport.Close()
				continue
			}
		}

		return nil
	}

	return fmt.Errorf("max reconnection attempts (%d) exceeded", r.config.MaxAttempts)
}
