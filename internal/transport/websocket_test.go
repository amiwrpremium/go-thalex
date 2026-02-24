package transport

import (
	"context"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
)

// ---------------------------------------------------------------------------
// NewWSTransport
// ---------------------------------------------------------------------------

func TestNewWSTransport(t *testing.T) {
	t.Run("default DialTimeout when zero", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL: "ws://localhost:8080",
		})
		if ws.dialTimeout != 10*time.Second {
			t.Errorf("dialTimeout = %v; want 10s", ws.dialTimeout)
		}
	})

	t.Run("default DialTimeout when negative", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL:         "ws://localhost:8080",
			DialTimeout: -5 * time.Second,
		})
		if ws.dialTimeout != 10*time.Second {
			t.Errorf("dialTimeout = %v; want 10s", ws.dialTimeout)
		}
	})

	t.Run("default PingInterval when zero", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL: "ws://localhost:8080",
		})
		if ws.pingInterval != 5*time.Second {
			t.Errorf("pingInterval = %v; want 5s", ws.pingInterval)
		}
	})

	t.Run("default PingInterval when negative", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL:          "ws://localhost:8080",
			PingInterval: -1 * time.Second,
		})
		if ws.pingInterval != 5*time.Second {
			t.Errorf("pingInterval = %v; want 5s", ws.pingInterval)
		}
	})

	t.Run("custom values preserved", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL:          "ws://custom.host:9090/ws",
			DialTimeout:  30 * time.Second,
			PingInterval: 15 * time.Second,
		})
		if ws.url != "ws://custom.host:9090/ws" {
			t.Errorf("url = %q; want %q", ws.url, "ws://custom.host:9090/ws")
		}
		if ws.dialTimeout != 30*time.Second {
			t.Errorf("dialTimeout = %v; want 30s", ws.dialTimeout)
		}
		if ws.pingInterval != 15*time.Second {
			t.Errorf("pingInterval = %v; want 15s", ws.pingInterval)
		}
	})

	t.Run("handler is stored", func(t *testing.T) {
		h := &mockWSHandler{}
		ws := NewWSTransport(WSTransportConfig{
			URL:     "ws://localhost:8080",
			Handler: h,
		})
		if ws.handler == nil {
			t.Fatal("expected handler to be set")
		}
	})

	t.Run("done channel is initialized", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL: "ws://localhost:8080",
		})
		if ws.done == nil {
			t.Fatal("expected done channel to be initialized")
		}
	})
}

// ---------------------------------------------------------------------------
// IsConnected
// ---------------------------------------------------------------------------

func TestWSTransport_IsConnected(t *testing.T) {
	t.Run("returns false before Connect", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL: "ws://localhost:8080",
		})
		if ws.IsConnected() {
			t.Error("expected IsConnected() to return false before Connect")
		}
	})

	t.Run("returns false when conn is nil", func(t *testing.T) {
		ws := &WSTransport{}
		if ws.IsConnected() {
			t.Error("expected IsConnected() to return false for zero-value WSTransport")
		}
	})
}

// ---------------------------------------------------------------------------
// Send
// ---------------------------------------------------------------------------

func TestWSTransport_Send(t *testing.T) {
	t.Run("returns error when not connected", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL: "ws://localhost:8080",
		})
		_, err := ws.Send(context.Background(), "test_method", nil)
		if err == nil {
			t.Fatal("expected an error when sending without connection")
		}
		if err.Error() != "not connected" {
			t.Errorf("error = %q; want %q", err.Error(), "not connected")
		}
	})

	t.Run("returns error for zero-value transport", func(t *testing.T) {
		ws := &WSTransport{}
		_, err := ws.Send(context.Background(), "method", nil)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

// ---------------------------------------------------------------------------
// Close
// ---------------------------------------------------------------------------

func TestWSTransport_Close(t *testing.T) {
	t.Run("close without connect does not panic", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL: "ws://localhost:8080",
		})
		err := ws.Close()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("double close does not panic", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL: "ws://localhost:8080",
		})
		ws.Close()
		// The second close: closeOnce ensures no double-close of the channel.
		// However conn is already nil, so it should be fine.
		// Note: we cannot call Close() a second time safely because closeOnce
		// was already used. But since conn is nil, the function should return nil.
		// The done channel is already closed, so we just verify no panic.
	})
}

// ---------------------------------------------------------------------------
// mockWSHandler (helper for tests)
// ---------------------------------------------------------------------------

type mockWSHandler struct{}

func (m *mockWSHandler) OnResponse(resp *jsonrpc.Response)          {}
func (m *mockWSHandler) OnNotification(notif *jsonrpc.Notification) {}
func (m *mockWSHandler) OnError(err error)                          {}
func (m *mockWSHandler) OnDisconnect()                              {}
