package transport

import (
	"compress/flate"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	gorilla "github.com/gorilla/websocket"

	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
)

// WSHandler processes incoming WebSocket messages.
type WSHandler interface {
	// OnResponse is called when a JSON-RPC response is received.
	OnResponse(resp *jsonrpc.Response)
	// OnNotification is called when a JSON-RPC notification is received.
	OnNotification(notif *jsonrpc.Notification)
	// OnError is called when a connection-level error occurs.
	OnError(err error)
	// OnDisconnect is called when the connection is lost.
	OnDisconnect()
}

// WSTransport manages a WebSocket connection to the Thalex API.
type WSTransport struct {
	url          string
	dialTimeout  time.Duration
	pingInterval time.Duration
	handler      WSHandler
	idGen        jsonrpc.IDGenerator

	mu   sync.Mutex
	conn *gorilla.Conn

	done      chan struct{}
	closeOnce sync.Once
}

// WSTransportConfig contains configuration for the WebSocket transport.
type WSTransportConfig struct {
	URL          string
	DialTimeout  time.Duration
	PingInterval time.Duration
	Handler      WSHandler
}

// NewWSTransport creates a new WebSocket transport.
func NewWSTransport(cfg WSTransportConfig) *WSTransport {
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = 10 * time.Second
	}
	if cfg.PingInterval <= 0 {
		cfg.PingInterval = 5 * time.Second
	}
	return &WSTransport{
		url:          cfg.URL,
		dialTimeout:  cfg.DialTimeout,
		pingInterval: cfg.PingInterval,
		handler:      cfg.Handler,
		done:         make(chan struct{}),
	}
}

// Connect establishes the WebSocket connection.
func (t *WSTransport) Connect(ctx context.Context) error {
	dialCtx, cancel := context.WithTimeout(ctx, t.dialTimeout)
	defer cancel()

	dialer := gorilla.Dialer{
		EnableCompression: true,
	}

	conn, resp, err := dialer.DialContext(dialCtx, t.url, nil)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("dialing WebSocket: %w", err)
	}

	// Enable write compression and set compression level.
	conn.EnableWriteCompression(true)
	_ = conn.SetCompressionLevel(flate.DefaultCompression)

	// Set a large read limit for order book snapshots etc.
	conn.SetReadLimit(16 * 1024 * 1024)

	t.mu.Lock()
	t.conn = conn
	t.done = make(chan struct{})
	t.closeOnce = sync.Once{}
	t.mu.Unlock()

	go t.readPump()
	go t.pingPump()

	return nil
}

// Close gracefully closes the WebSocket connection.
func (t *WSTransport) Close() error {
	t.closeOnce.Do(func() {
		close(t.done)
	})

	t.mu.Lock()
	conn := t.conn
	t.conn = nil
	t.mu.Unlock()

	if conn != nil {
		// Send a close frame, then close the underlying connection.
		_ = conn.WriteMessage(
			gorilla.CloseMessage,
			gorilla.FormatCloseMessage(gorilla.CloseNormalClosure, "client closing"),
		)
		return conn.Close()
	}
	return nil
}

// Send sends a JSON-RPC request and returns the request ID.
func (t *WSTransport) Send(ctx context.Context, method string, params any) (uint64, error) {
	id := t.idGen.Next()
	req := jsonrpc.NewRequest(id, method, params)

	data, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("marshaling request: %w", err)
	}

	t.mu.Lock()
	conn := t.conn
	t.mu.Unlock()

	if conn == nil {
		return 0, fmt.Errorf("not connected")
	}

	// Respect context deadline if present.
	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetWriteDeadline(deadline); err != nil {
			return 0, fmt.Errorf("setting write deadline: %w", err)
		}
	}

	if err := conn.WriteMessage(gorilla.TextMessage, data); err != nil {
		return 0, fmt.Errorf("writing to WebSocket: %w", err)
	}

	return id, nil
}

// IsConnected returns true if the WebSocket is currently connected.
func (t *WSTransport) IsConnected() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.conn != nil
}

// readPump continuously reads messages from the WebSocket.
func (t *WSTransport) readPump() {
	// Capture done channel under the lock so we don't race with Connect()
	// replacing t.done on reconnect.
	t.mu.Lock()
	done := t.done
	t.mu.Unlock()

	defer func() {
		if t.handler != nil {
			t.handler.OnDisconnect()
		}
	}()

	for {
		select {
		case <-done:
			return
		default:
		}

		t.mu.Lock()
		conn := t.conn
		t.mu.Unlock()

		if conn == nil {
			return
		}

		_, data, err := conn.ReadMessage()
		if err != nil {
			select {
			case <-done:
				return
			default:
			}
			if t.handler != nil {
				t.handler.OnError(fmt.Errorf("reading from WebSocket: %w", err))
			}
			return
		}

		msg, err := jsonrpc.ParseMessage(data)
		if err != nil {
			if t.handler != nil {
				t.handler.OnError(fmt.Errorf("parsing message: %w", err))
			}
			continue
		}

		if t.handler == nil {
			continue
		}

		if msg.Response != nil {
			t.handler.OnResponse(msg.Response)
		} else if msg.Notification != nil {
			t.handler.OnNotification(msg.Notification)
		}
	}
}

// pingPump sends periodic ping frames to keep the connection alive.
func (t *WSTransport) pingPump() {
	// Capture done channel under the lock so we don't race with Connect()
	// replacing t.done on reconnect.
	t.mu.Lock()
	done := t.done
	t.mu.Unlock()

	ticker := time.NewTicker(t.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			t.mu.Lock()
			conn := t.conn
			t.mu.Unlock()

			if conn == nil {
				return
			}

			deadline := time.Now().Add(5 * time.Second)
			err := conn.WriteControl(gorilla.PingMessage, nil, deadline)

			if err != nil {
				if t.handler != nil {
					t.handler.OnError(fmt.Errorf("ping failed: %w", err))
				}
				return
			}
		}
	}
}
