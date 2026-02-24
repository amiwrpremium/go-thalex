package ws

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/amiwrpremium/go-thalex/apierr"
	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
	"github.com/amiwrpremium/go-thalex/internal/transport"
	"github.com/amiwrpremium/go-thalex/types"
)

type pendingCall struct {
	result chan *jsonrpc.Response
}

// Client provides access to the Thalex WebSocket JSON-RPC API.
type Client struct {
	transport   *transport.WSTransport
	reconnector *transport.Reconnector
	cfg         config.ClientConfig

	mu      sync.Mutex
	pending map[uint64]*pendingCall

	subMu    sync.RWMutex
	handlers map[string]any

	onError func(error)
}

// NewClient creates a new WebSocket API client.
func NewClient(opts ...config.ClientOption) *Client {
	cfg := config.DefaultClientConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	ws := &Client{
		cfg:      cfg,
		pending:  make(map[uint64]*pendingCall),
		handlers: make(map[string]any),
	}
	ws.transport = transport.NewWSTransport(transport.WSTransportConfig{
		URL:          cfg.Network.WebSocketURL(),
		DialTimeout:  cfg.WSDialTimeout,
		PingInterval: cfg.WSPingInterval,
		Handler:      ws,
	})
	if cfg.WSReconnect {
		ws.reconnector = transport.NewReconnector(ws.transport, transport.ReconnectConfig{
			Enabled:     true,
			MaxAttempts: cfg.WSMaxReconnects,
			BaseWait:    cfg.WSReconnectWait,
			OnReconnect: ws.onReconnect,
		})
	}
	return ws
}

// Connect establishes the WebSocket connection.
func (ws *Client) Connect(ctx context.Context) error {
	if err := ws.transport.Connect(ctx); err != nil {
		return err
	}
	if ws.reconnector != nil {
		ws.reconnector.Start(ctx)
	}
	return nil
}

// Close gracefully closes the WebSocket connection.
func (ws *Client) Close() error {
	if ws.reconnector != nil {
		ws.reconnector.Stop()
	}
	ws.mu.Lock()
	for id, pc := range ws.pending {
		close(pc.result)
		delete(ws.pending, id)
	}
	ws.mu.Unlock()
	return ws.transport.Close()
}

// OnErrorHandler registers a callback for connection-level errors.
func (ws *Client) OnErrorHandler(fn func(error)) {
	ws.onError = fn
}

// IsConnected returns true if the WebSocket is currently connected.
func (ws *Client) IsConnected() bool {
	return ws.transport.IsConnected()
}

// call sends a JSON-RPC request and waits for the response.
func (ws *Client) call(ctx context.Context, method string, params any, result any) error {
	pc := &pendingCall{result: make(chan *jsonrpc.Response, 1)}

	id, err := ws.transport.Send(ctx, method, params)
	if err != nil {
		return err
	}

	ws.mu.Lock()
	ws.pending[id] = pc
	ws.mu.Unlock()

	defer func() {
		ws.mu.Lock()
		delete(ws.pending, id)
		ws.mu.Unlock()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case resp, ok := <-pc.result:
		if !ok {
			return &apierr.ConnectionError{Message: "connection closed while waiting for response"}
		}
		if resp.Error != nil {
			return &apierr.APIError{Code: resp.Error.Code, Message: resp.Error.Message}
		}
		if result != nil && resp.Result != nil {
			return json.Unmarshal(resp.Result, result)
		}
		return nil
	}
}

// callNoResult sends a JSON-RPC request expecting a null result.
func (ws *Client) callNoResult(ctx context.Context, method string, params any) error {
	return ws.call(ctx, method, params, nil)
}

// --- WSHandler interface implementation ---

// OnResponse dispatches a JSON-RPC response to the pending call.
func (ws *Client) OnResponse(resp *jsonrpc.Response) {
	if resp.ID == nil {
		return
	}
	ws.mu.Lock()
	pc, ok := ws.pending[*resp.ID]
	ws.mu.Unlock()
	if ok {
		pc.result <- resp
	}
}

// OnNotification dispatches a JSON-RPC notification to the subscription handler.
func (ws *Client) OnNotification(notif *jsonrpc.Notification) {
	ws.subMu.RLock()
	handler, ok := ws.handlers[notif.Method]
	ws.subMu.RUnlock()
	if ok {
		ws.dispatchNotification(handler, notif.Params)
	}
}

// OnError handles connection-level errors.
func (ws *Client) OnError(err error) {
	if ws.onError != nil {
		ws.onError(err)
	}
}

// OnDisconnect handles connection loss and triggers reconnection.
func (ws *Client) OnDisconnect() {
	if ws.reconnector != nil {
		go func() {
			_ = ws.reconnector.TriggerReconnect(context.Background())
		}()
	}
}

func (ws *Client) onReconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if ws.cfg.Credentials != nil {
		if err := ws.Login(ctx); err != nil {
			return err
		}
	}

	ws.subMu.RLock()
	var pub, priv []string
	for ch := range ws.handlers {
		if isPrivateChannel(ch) {
			priv = append(priv, ch)
		} else {
			pub = append(pub, ch)
		}
	}
	ws.subMu.RUnlock()

	if len(pub) > 0 {
		_ = ws.callNoResult(ctx, "public/subscribe", map[string]any{"channels": pub})
	}
	if len(priv) > 0 {
		_ = ws.callNoResult(ctx, "private/subscribe", map[string]any{"channels": priv})
	}
	return nil
}

func isPrivateChannel(ch string) bool {
	return strings.HasPrefix(ch, "account.") || strings.HasPrefix(ch, "session.") ||
		strings.HasPrefix(ch, "user.") || strings.HasPrefix(ch, "mm.")
}

func (ws *Client) dispatchNotification(handler any, data json.RawMessage) {
	go func() {
		switch fn := handler.(type) {
		case func(types.BookUpdate):
			var v types.BookUpdate
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func(types.Ticker):
			var v types.Ticker
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func(types.LightweightTicker):
			var v types.LightweightTicker
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.RecentTrade):
			var v []types.RecentTrade
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func(types.IndexPrice):
			var v types.IndexPrice
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.Instrument):
			var v []types.Instrument
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.OrderStatus):
			var v []types.OrderStatus
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.PortfolioEntry):
			var v []types.PortfolioEntry
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func(types.AccountSummary):
			var v types.AccountSummary
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.Trade):
			var v []types.Trade
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.OrderHistory):
			var v []types.OrderHistory
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.ConditionalOrder):
			var v []types.ConditionalOrder
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.Bot):
			var v []types.Bot
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.Rfq):
			var v []types.Rfq
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.RfqOrder):
			var v []types.RfqOrder
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func(types.MMProtectionUpdate):
			var v types.MMProtectionUpdate
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func(types.Notification):
			var v types.Notification
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func(types.SystemEvent):
			var v types.SystemEvent
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func([]types.Banner):
			var v []types.Banner
			if json.Unmarshal(data, &v) == nil {
				fn(v)
			}
		case func(json.RawMessage):
			fn(data)
		}
	}()
}
