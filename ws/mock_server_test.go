package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gorilla "github.com/gorilla/websocket"

	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
	"github.com/amiwrpremium/go-thalex/internal/transport"
)

// rpcHandler is a function that receives a parsed JSON-RPC request and returns
// either a result payload or a JSON-RPC error.
type rpcHandler func(req *jsonrpc.Request) (result json.RawMessage, rpcErr *jsonrpc.Error)

// newMockWSServer creates an httptest.Server that upgrades to WebSocket and
// responds to JSON-RPC requests using the provided handler function.
func newMockWSServer(t *testing.T, handler rpcHandler) *httptest.Server {
	t.Helper()
	upgrader := gorilla.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("mock server upgrade error: %v", err)
			return
		}
		defer conn.Close()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return // client closed
			}

			var req jsonrpc.Request
			if err := json.Unmarshal(data, &req); err != nil {
				t.Logf("mock server: failed to parse request: %v", err)
				continue
			}

			result, rpcErr := handler(&req)

			resp := struct {
				JSONRPC string          `json:"jsonrpc"`
				ID      uint64          `json:"id"`
				Result  json.RawMessage `json:"result,omitempty"`
				Error   *jsonrpc.Error  `json:"error,omitempty"`
			}{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  result,
				Error:   rpcErr,
			}

			respData, _ := json.Marshal(resp)
			if err := conn.WriteMessage(gorilla.TextMessage, respData); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

// wsURLFromHTTP converts an httptest server URL (http://...) to a ws:// URL.
func wsURLFromHTTP(httpURL string) string {
	return "ws" + strings.TrimPrefix(httpURL, "http")
}

// newConnectedClient creates a Client connected to a mock WS server that
// uses the given handler for all requests. The client is connected and ready.
func newConnectedClient(t *testing.T, handler rpcHandler) *Client {
	t.Helper()

	srv := newMockWSServer(t, handler)

	cfg := config.DefaultClientConfig()
	c := &Client{
		cfg:      cfg,
		pending:  make(map[uint64]*pendingCall),
		handlers: make(map[string]any),
	}
	c.transport = transport.NewWSTransport(transport.WSTransportConfig{
		URL:          wsURLFromHTTP(srv.URL),
		DialTimeout:  cfg.WSDialTimeout,
		PingInterval: 60 * time.Second, // slow pings so they don't interfere
		Handler:      c,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.transport.Connect(ctx); err != nil {
		t.Fatalf("failed to connect to mock server: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })

	return c
}

// echoNull is a simple rpcHandler that returns JSON null for every request.
func echoNull(_ *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
	return json.RawMessage(`null`), nil
}

// methodRouter builds an rpcHandler from a map of method->handler.
// Unmatched methods return a JSON-RPC error.
func methodRouter(routes map[string]rpcHandler) rpcHandler {
	return func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if h, ok := routes[req.Method]; ok {
			return h(req)
		}
		return nil, &jsonrpc.Error{Code: -32601, Message: "method not found: " + req.Method}
	}
}
