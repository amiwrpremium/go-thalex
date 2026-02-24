package transport

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// NewReconnector
// ---------------------------------------------------------------------------

func TestNewReconnector(t *testing.T) {
	t.Run("default BaseWait when zero", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{})
		if r.config.BaseWait != 1*time.Second {
			t.Errorf("BaseWait = %v; want 1s", r.config.BaseWait)
		}
	})

	t.Run("default BaseWait when negative", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{BaseWait: -1 * time.Second})
		if r.config.BaseWait != 1*time.Second {
			t.Errorf("BaseWait = %v; want 1s", r.config.BaseWait)
		}
	})

	t.Run("default MaxWait when zero", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{})
		if r.config.MaxWait != 30*time.Second {
			t.Errorf("MaxWait = %v; want 30s", r.config.MaxWait)
		}
	})

	t.Run("default MaxWait when negative", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{MaxWait: -5 * time.Second})
		if r.config.MaxWait != 30*time.Second {
			t.Errorf("MaxWait = %v; want 30s", r.config.MaxWait)
		}
	})

	t.Run("custom values preserved", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{
			Enabled:     true,
			MaxAttempts: 10,
			BaseWait:    2 * time.Second,
			MaxWait:     60 * time.Second,
		})
		if !r.config.Enabled {
			t.Error("Enabled should be true")
		}
		if r.config.MaxAttempts != 10 {
			t.Errorf("MaxAttempts = %d; want 10", r.config.MaxAttempts)
		}
		if r.config.BaseWait != 2*time.Second {
			t.Errorf("BaseWait = %v; want 2s", r.config.BaseWait)
		}
		if r.config.MaxWait != 60*time.Second {
			t.Errorf("MaxWait = %v; want 60s", r.config.MaxWait)
		}
	})

	t.Run("transport is stored", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{})
		if r.transport != ws {
			t.Error("expected transport to be stored")
		}
	})

	t.Run("OnReconnect callback is stored", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		called := false
		r := NewReconnector(ws, ReconnectConfig{
			OnReconnect: func() error {
				called = true
				return nil
			},
		})
		if r.config.OnReconnect == nil {
			t.Fatal("expected OnReconnect to be set")
		}
		// Invoke to verify it's the right function.
		r.config.OnReconnect()
		if !called {
			t.Error("expected OnReconnect callback to be invoked")
		}
	})
}

// ---------------------------------------------------------------------------
// Start
// ---------------------------------------------------------------------------

func TestReconnector_Start(t *testing.T) {
	t.Run("does nothing when Enabled is false", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{Enabled: false})
		r.Start(context.Background())

		// active should remain false.
		r.mu.Lock()
		active := r.active
		r.mu.Unlock()
		if active {
			t.Error("expected active to be false when Enabled is false")
		}
	})

	t.Run("sets active to true when Enabled", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{Enabled: true})
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r.Start(ctx)
		defer r.Stop()

		r.mu.Lock()
		active := r.active
		r.mu.Unlock()
		if !active {
			t.Error("expected active to be true after Start")
		}
	})

	t.Run("does not start twice - idempotent", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{Enabled: true})
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r.Start(ctx)
		defer r.Stop()

		// Capture the cancelFunc after first Start.
		r.mu.Lock()
		firstCancel := r.cancelFunc
		r.mu.Unlock()

		// Start again should be a no-op.
		r.Start(ctx)

		r.mu.Lock()
		secondCancel := r.cancelFunc
		r.mu.Unlock()

		// The cancelFunc pointer should be the same, indicating no restart.
		if fmt.Sprintf("%p", firstCancel) != fmt.Sprintf("%p", secondCancel) {
			t.Error("expected Start to be idempotent; cancelFunc changed")
		}
	})

	t.Run("sets cancelFunc", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{Enabled: true})
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r.Start(ctx)
		defer r.Stop()

		r.mu.Lock()
		hasCancel := r.cancelFunc != nil
		r.mu.Unlock()
		if !hasCancel {
			t.Error("expected cancelFunc to be set after Start")
		}
	})
}

// ---------------------------------------------------------------------------
// Stop
// ---------------------------------------------------------------------------

func TestReconnector_Stop(t *testing.T) {
	t.Run("sets active to false", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{Enabled: true})
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r.Start(ctx)

		r.mu.Lock()
		activeBefore := r.active
		r.mu.Unlock()
		if !activeBefore {
			t.Fatal("expected active to be true before Stop")
		}

		r.Stop()

		r.mu.Lock()
		activeAfter := r.active
		r.mu.Unlock()
		if activeAfter {
			t.Error("expected active to be false after Stop")
		}
	})

	t.Run("clears cancelFunc", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{Enabled: true})
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r.Start(ctx)
		r.Stop()

		r.mu.Lock()
		cf := r.cancelFunc
		r.mu.Unlock()
		if cf != nil {
			t.Error("expected cancelFunc to be nil after Stop")
		}
	})

	t.Run("handles nil cancelFunc gracefully", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{})
		// cancelFunc is nil by default; Stop should not panic.
		r.Stop()

		r.mu.Lock()
		active := r.active
		r.mu.Unlock()
		if active {
			t.Error("expected active to be false")
		}
	})

	t.Run("double stop does not panic", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{Enabled: true})
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r.Start(ctx)
		r.Stop()
		r.Stop() // Should not panic.
	})
}

// ---------------------------------------------------------------------------
// TriggerReconnect
// ---------------------------------------------------------------------------

func TestReconnector_TriggerReconnect(t *testing.T) {
	t.Run("returns error when disabled", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{Enabled: false})

		err := r.TriggerReconnect(context.Background())
		if err == nil {
			t.Fatal("expected an error when reconnection is disabled")
		}
		if err.Error() != "reconnection is disabled" {
			t.Errorf("error = %q; want %q", err.Error(), "reconnection is disabled")
		}
	})

	t.Run("with cancelled context returns context error", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
		r := NewReconnector(ws, ReconnectConfig{
			Enabled:     true,
			MaxAttempts: 1,
			BaseWait:    10 * time.Millisecond,
		})

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately.

		err := r.TriggerReconnect(ctx)
		if err == nil {
			t.Fatal("expected an error")
		}
		// Should get context.Canceled since ctx is already cancelled.
		if err != context.Canceled {
			t.Errorf("error = %v; want context.Canceled", err)
		}
	})

	t.Run("max attempts exceeded returns error", func(t *testing.T) {
		// WSTransport will fail to connect because there's no real WS server.
		ws := NewWSTransport(WSTransportConfig{
			URL:         "ws://127.0.0.1:1", // Nothing listening here.
			DialTimeout: 50 * time.Millisecond,
		})
		r := NewReconnector(ws, ReconnectConfig{
			Enabled:     true,
			MaxAttempts: 2,
			BaseWait:    time.Millisecond,
			MaxWait:     time.Millisecond,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := r.TriggerReconnect(ctx)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "max reconnection attempts (2) exceeded" {
			t.Errorf("error = %q; want max reconnection attempts exceeded", err.Error())
		}
	})
}
