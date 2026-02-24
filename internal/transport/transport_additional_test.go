package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	gorilla "github.com/gorilla/websocket"

	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// wsTestServer creates a test HTTP server that upgrades to WebSocket.
// The handler function receives the connection for custom behavior.
// Returns the server and its "ws://" URL.
func wsTestServer(t *testing.T, handler func(conn *gorilla.Conn)) (*httptest.Server, string) {
	t.Helper()
	upgrader := gorilla.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("upgrade error: %v", err)
			return
		}
		defer conn.Close()
		handler(conn)
	}))
	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	return server, wsURL
}

// recordingHandler records calls for assertions.
type recordingHandler struct {
	mu            sync.Mutex
	responses     []*jsonrpc.Response
	notifications []*jsonrpc.Notification
	errors        []error
	disconnects   int32
	disconnectCh  chan struct{}
}

func newRecordingHandler() *recordingHandler {
	return &recordingHandler{
		disconnectCh: make(chan struct{}, 1),
	}
}

func (h *recordingHandler) OnResponse(resp *jsonrpc.Response) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.responses = append(h.responses, resp)
}

func (h *recordingHandler) OnNotification(notif *jsonrpc.Notification) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.notifications = append(h.notifications, notif)
}

func (h *recordingHandler) OnError(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.errors = append(h.errors, err)
}

func (h *recordingHandler) OnDisconnect() {
	atomic.AddInt32(&h.disconnects, 1)
	select {
	case h.disconnectCh <- struct{}{}:
	default:
	}
}

// ---------------------------------------------------------------------------
// WSTransport – Connect and Send with real WebSocket server
// ---------------------------------------------------------------------------

func TestWSTransport_ConnectAndSend(t *testing.T) {
	t.Run("connect and send successfully", func(t *testing.T) {
		var receivedMsg []byte
		serverDone := make(chan struct{})
		server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			receivedMsg = msg
			// Send a response back.
			var req jsonrpc.Request
			json.Unmarshal(msg, &req)
			resp := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"result":  "ok",
			}
			conn.WriteJSON(resp)
			close(serverDone)
			// Keep the connection alive briefly for readPump.
			time.Sleep(100 * time.Millisecond)
		})
		defer server.Close()

		handler := newRecordingHandler()
		ws := NewWSTransport(WSTransportConfig{
			URL:          wsURL,
			DialTimeout:  5 * time.Second,
			PingInterval: 1 * time.Hour, // Disable pings for this test.
			Handler:      handler,
		})

		err := ws.Connect(context.Background())
		if err != nil {
			t.Fatalf("Connect failed: %v", err)
		}
		defer ws.Close()

		if !ws.IsConnected() {
			t.Error("expected IsConnected() to be true after Connect")
		}

		id, err := ws.Send(context.Background(), "test.method", map[string]string{"key": "val"})
		if err != nil {
			t.Fatalf("Send failed: %v", err)
		}
		if id == 0 {
			t.Error("expected non-zero request ID")
		}

		// Wait for server to receive the message.
		select {
		case <-serverDone:
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for server to receive message")
		}

		if receivedMsg == nil {
			t.Fatal("server did not receive a message")
		}

		var parsed jsonrpc.Request
		if err := json.Unmarshal(receivedMsg, &parsed); err != nil {
			t.Fatalf("failed to parse received message: %v", err)
		}
		if parsed.Method != "test.method" {
			t.Errorf("Method = %q; want %q", parsed.Method, "test.method")
		}
		if parsed.ID != id {
			t.Errorf("ID = %d; want %d", parsed.ID, id)
		}

		// Give readPump time to process the response.
		time.Sleep(50 * time.Millisecond)

		handler.mu.Lock()
		respCount := len(handler.responses)
		handler.mu.Unlock()
		if respCount == 0 {
			t.Error("expected at least one response to be received by handler")
		}
	})

	t.Run("connect fails with invalid URL", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL:         "ws://127.0.0.1:1/invalid",
			DialTimeout: 100 * time.Millisecond,
		})
		err := ws.Connect(context.Background())
		if err == nil {
			t.Fatal("expected Connect to fail with invalid URL")
		}
		if !strings.Contains(err.Error(), "dialing WebSocket") {
			t.Errorf("error = %q; want to contain 'dialing WebSocket'", err.Error())
		}
	})

	t.Run("connect fails when context is cancelled", func(t *testing.T) {
		ws := NewWSTransport(WSTransportConfig{
			URL:         "ws://127.0.0.1:1/invalid",
			DialTimeout: 5 * time.Second,
		})
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := ws.Connect(ctx)
		if err == nil {
			t.Fatal("expected Connect to fail with cancelled context")
		}
	})

	t.Run("send with deadline context", func(t *testing.T) {
		serverReady := make(chan struct{})
		server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
			close(serverReady)
			// Read and discard.
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					return
				}
			}
		})
		defer server.Close()

		ws := NewWSTransport(WSTransportConfig{
			URL:          wsURL,
			PingInterval: 1 * time.Hour,
		})
		err := ws.Connect(context.Background())
		if err != nil {
			t.Fatalf("Connect failed: %v", err)
		}
		defer ws.Close()

		<-serverReady

		// Send with a deadline context.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id, err := ws.Send(ctx, "deadline.method", nil)
		if err != nil {
			t.Fatalf("Send with deadline failed: %v", err)
		}
		if id == 0 {
			t.Error("expected non-zero ID")
		}
	})
}

// ---------------------------------------------------------------------------
// WSTransport – readPump dispatching
// ---------------------------------------------------------------------------

func TestWSTransport_ReadPump_Notification(t *testing.T) {
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		// Send a notification.
		notif := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "subscription",
			"params":  map[string]string{"channel": "trades"},
		}
		conn.WriteJSON(notif)
		// Keep alive briefly.
		time.Sleep(200 * time.Millisecond)
	})
	defer server.Close()

	handler := newRecordingHandler()
	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		PingInterval: 1 * time.Hour,
		Handler:      handler,
	})

	err := ws.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer ws.Close()

	// Wait for the notification to arrive.
	time.Sleep(100 * time.Millisecond)

	handler.mu.Lock()
	notifCount := len(handler.notifications)
	handler.mu.Unlock()
	if notifCount == 0 {
		t.Error("expected at least one notification to be received")
	}
}

func TestWSTransport_ReadPump_InvalidJSON(t *testing.T) {
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		// Send invalid JSON.
		conn.WriteMessage(gorilla.TextMessage, []byte(`{not valid json}`))
		// Then send a valid response so the readPump continues.
		resp := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  "ok",
		}
		conn.WriteJSON(resp)
		time.Sleep(200 * time.Millisecond)
	})
	defer server.Close()

	handler := newRecordingHandler()
	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		PingInterval: 1 * time.Hour,
		Handler:      handler,
	})

	err := ws.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer ws.Close()

	// Wait for processing.
	time.Sleep(100 * time.Millisecond)

	handler.mu.Lock()
	errCount := len(handler.errors)
	handler.mu.Unlock()
	if errCount == 0 {
		t.Error("expected OnError to be called for invalid JSON")
	}
}

func TestWSTransport_ReadPump_NilHandler(t *testing.T) {
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		// Send a message that would normally trigger handler calls.
		resp := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  "ok",
		}
		conn.WriteJSON(resp)
		time.Sleep(200 * time.Millisecond)
	})
	defer server.Close()

	// No handler set -- should not panic.
	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		PingInterval: 1 * time.Hour,
		Handler:      nil,
	})

	err := ws.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Give it time to read the message without panicking.
	time.Sleep(100 * time.Millisecond)

	ws.Close()
}

// ---------------------------------------------------------------------------
// WSTransport – readPump disconnect
// ---------------------------------------------------------------------------

func TestWSTransport_ReadPump_OnDisconnect(t *testing.T) {
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		// Close immediately to trigger disconnect.
		conn.Close()
	})
	defer server.Close()

	handler := newRecordingHandler()
	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		PingInterval: 1 * time.Hour,
		Handler:      handler,
	})

	err := ws.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer ws.Close()

	// Wait for OnDisconnect.
	select {
	case <-handler.disconnectCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnDisconnect")
	}

	if atomic.LoadInt32(&handler.disconnects) == 0 {
		t.Error("expected OnDisconnect to be called")
	}
}

// ---------------------------------------------------------------------------
// WSTransport – Close sets conn to nil
// ---------------------------------------------------------------------------

func TestWSTransport_CloseAfterConnect(t *testing.T) {
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		// Keep alive until test closes.
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		PingInterval: 1 * time.Hour,
		Handler:      newRecordingHandler(),
	})

	err := ws.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if !ws.IsConnected() {
		t.Error("expected IsConnected() to be true")
	}

	err = ws.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if ws.IsConnected() {
		t.Error("expected IsConnected() to be false after Close")
	}
}

// ---------------------------------------------------------------------------
// WSTransport – pingPump
// ---------------------------------------------------------------------------

func TestWSTransport_PingPump(t *testing.T) {
	var pingReceived atomic.Int32
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		conn.SetPingHandler(func(appData string) error {
			pingReceived.Add(1)
			return conn.WriteControl(gorilla.PongMessage, []byte(appData), time.Now().Add(time.Second))
		})
		// Read messages to keep the connection alive.
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		PingInterval: 50 * time.Millisecond,
		Handler:      newRecordingHandler(),
	})

	err := ws.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Wait long enough for at least one ping.
	time.Sleep(200 * time.Millisecond)

	ws.Close()

	if pingReceived.Load() == 0 {
		t.Error("expected at least one ping to be received by server")
	}
}

// ---------------------------------------------------------------------------
// HTTPTransport – DoPrivateGET no query params
// ---------------------------------------------------------------------------

func TestDoPrivateGET_NoQueryParams(t *testing.T) {
	var receivedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.String()
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
	}))
	defer server.Close()

	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       server.URL,
		RetryBaseWait: time.Millisecond,
	})
	err := tr.DoPrivateGET(context.Background(), "/balances", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(receivedURL, "?") {
		t.Errorf("URL %q should not contain '?' with no params", receivedURL)
	}
}

func TestDoPrivateGET_EmptyQueryParams(t *testing.T) {
	var receivedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.String()
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
	}))
	defer server.Close()

	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       server.URL,
		RetryBaseWait: time.Millisecond,
	})
	err := tr.DoPrivateGET(context.Background(), "/info", url.Values{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// url.Values{} has length 0, so no "?" should be appended.
	if strings.Contains(receivedURL, "?") {
		t.Errorf("URL %q should not contain '?' with empty params", receivedURL)
	}
}

// ---------------------------------------------------------------------------
// HTTPTransport – DoPrivatePOST body marshal error
// ---------------------------------------------------------------------------

func TestDoPrivatePOST_MarshalError(t *testing.T) {
	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       "http://localhost",
		RetryBaseWait: time.Millisecond,
	})

	// Channels cannot be marshaled to JSON.
	body := map[string]interface{}{"bad": make(chan int)}
	err := tr.DoPrivatePOST(context.Background(), "/test", body, nil)
	if err == nil {
		t.Fatal("expected marshal error")
	}
	if !strings.Contains(err.Error(), "marshaling request body") {
		t.Errorf("error = %q; want to contain 'marshaling request body'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// HTTPTransport – DoPrivatePOST result parsing
// ---------------------------------------------------------------------------

func TestDoPrivatePOST_SuccessfulResultParsing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(apiResponse{
			Result: json.RawMessage(`{"order_id":"ABC123"}`),
		})
	}))
	defer server.Close()

	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       server.URL,
		RetryBaseWait: time.Millisecond,
	})

	var result struct {
		OrderID string `json:"order_id"`
	}
	err := tr.DoPrivatePOST(context.Background(), "/order", map[string]string{"side": "buy"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ABC123" {
		t.Errorf("OrderID = %q; want %q", result.OrderID, "ABC123")
	}
}

// ---------------------------------------------------------------------------
// HTTPTransport – DoPrivatePOST error paths
// ---------------------------------------------------------------------------

func TestDoPrivatePOST_ServerError500Retries(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := attempts.Add(1)
		if count < 3 {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`"done"`)})
	}))
	defer server.Close()

	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       server.URL,
		MaxRetries:    5,
		RetryBaseWait: time.Millisecond,
	})
	var result string
	err := tr.DoPrivatePOST(context.Background(), "/retry", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "done" {
		t.Errorf("result = %q; want %q", result, "done")
	}
}

func TestDoPrivatePOST_ClientError400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(apiResponse{
			Error: &apiError{Code: 1001, Message: "Invalid"},
		})
	}))
	defer server.Close()

	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       server.URL,
		RetryBaseWait: time.Millisecond,
	})
	err := tr.DoPrivatePOST(context.Background(), "/bad", nil, nil)
	if err == nil {
		t.Fatal("expected an error")
	}
	apiErr, ok := err.(*apiError)
	if !ok {
		t.Fatalf("expected *apiError, got %T", err)
	}
	if apiErr.Code != 1001 {
		t.Errorf("Code = %d; want 1001", apiErr.Code)
	}
}

// ---------------------------------------------------------------------------
// HTTPTransport – setAuthHeaders edge cases
// ---------------------------------------------------------------------------

func TestSetAuthHeaders_NoTokenFunc_NoAccountNumber(t *testing.T) {
	var receivedAuth, receivedAccount string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		receivedAccount = r.Header.Get("X-Thalex-Account")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
	}))
	defer server.Close()

	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       server.URL,
		RetryBaseWait: time.Millisecond,
	})
	tr.DoPrivateGET(context.Background(), "/test", nil, nil)

	if receivedAuth != "" {
		t.Errorf("Authorization = %q; want empty", receivedAuth)
	}
	if receivedAccount != "" {
		t.Errorf("X-Thalex-Account = %q; want empty", receivedAccount)
	}
}

// ---------------------------------------------------------------------------
// HTTPTransport – doWithRetry context cancellation during retry wait
// ---------------------------------------------------------------------------

func TestDoWithRetry_ContextCancelDuringNetworkError(t *testing.T) {
	// Server that immediately closes connections.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       server.URL,
		MaxRetries:    100,
		RetryBaseWait: 100 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := tr.DoPublic(ctx, "/slow", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	// Should contain context error because context is cancelled during the retry wait.
	if !strings.Contains(err.Error(), "context") {
		t.Errorf("error = %q; want to contain 'context'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// Reconnector – reconnect with successful connection
// ---------------------------------------------------------------------------

func TestReconnector_TriggerReconnect_Success(t *testing.T) {
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	handler := newRecordingHandler()
	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		DialTimeout:  5 * time.Second,
		PingInterval: 1 * time.Hour,
		Handler:      handler,
	})

	callbackCalled := false
	r := NewReconnector(ws, ReconnectConfig{
		Enabled:     true,
		MaxAttempts: 3,
		BaseWait:    time.Millisecond,
		MaxWait:     10 * time.Millisecond,
		OnReconnect: func() error {
			callbackCalled = true
			return nil
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.TriggerReconnect(ctx)
	if err != nil {
		t.Fatalf("TriggerReconnect failed: %v", err)
	}
	defer ws.Close()

	if !callbackCalled {
		t.Error("expected OnReconnect callback to be called")
	}

	if !ws.IsConnected() {
		t.Error("expected transport to be connected after reconnect")
	}
}

// ---------------------------------------------------------------------------
// Reconnector – OnReconnect callback fails triggers retry
// ---------------------------------------------------------------------------

func TestReconnector_TriggerReconnect_OnReconnectFails(t *testing.T) {
	var connectCount atomic.Int32
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		connectCount.Add(1)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	handler := newRecordingHandler()
	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		DialTimeout:  5 * time.Second,
		PingInterval: 1 * time.Hour,
		Handler:      handler,
	})

	var callbackAttempts atomic.Int32
	r := NewReconnector(ws, ReconnectConfig{
		Enabled:     true,
		MaxAttempts: 3,
		BaseWait:    time.Millisecond,
		MaxWait:     5 * time.Millisecond,
		OnReconnect: func() error {
			count := callbackAttempts.Add(1)
			if count < 3 {
				return fmt.Errorf("callback failed attempt %d", count)
			}
			return nil
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := r.TriggerReconnect(ctx)
	if err != nil {
		t.Fatalf("TriggerReconnect failed: %v", err)
	}
	defer ws.Close()

	if callbackAttempts.Load() < 3 {
		t.Errorf("expected at least 3 callback attempts, got %d", callbackAttempts.Load())
	}
}

// ---------------------------------------------------------------------------
// Reconnector – max attempts exceeded with OnReconnect always failing
// ---------------------------------------------------------------------------

func TestReconnector_TriggerReconnect_MaxAttemptsWithCallbackAlwaysFailing(t *testing.T) {
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		DialTimeout:  5 * time.Second,
		PingInterval: 1 * time.Hour,
		Handler:      newRecordingHandler(),
	})

	r := NewReconnector(ws, ReconnectConfig{
		Enabled:     true,
		MaxAttempts: 2,
		BaseWait:    time.Millisecond,
		MaxWait:     5 * time.Millisecond,
		OnReconnect: func() error {
			return fmt.Errorf("always fails")
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := r.TriggerReconnect(ctx)
	if err == nil {
		t.Fatal("expected error from max attempts exceeded")
	}
	if !strings.Contains(err.Error(), "max reconnection attempts (2) exceeded") {
		t.Errorf("error = %q; want 'max reconnection attempts (2) exceeded'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// Reconnector – reconnect with nil OnReconnect callback
// ---------------------------------------------------------------------------

func TestReconnector_TriggerReconnect_NilOnReconnect(t *testing.T) {
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		DialTimeout:  5 * time.Second,
		PingInterval: 1 * time.Hour,
		Handler:      newRecordingHandler(),
	})

	r := NewReconnector(ws, ReconnectConfig{
		Enabled:     true,
		MaxAttempts: 3,
		BaseWait:    time.Millisecond,
		MaxWait:     5 * time.Millisecond,
		// OnReconnect is nil -- should succeed without callback.
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.TriggerReconnect(ctx)
	if err != nil {
		t.Fatalf("TriggerReconnect failed: %v", err)
	}
	defer ws.Close()

	if !ws.IsConnected() {
		t.Error("expected transport to be connected after reconnect")
	}
}

// ---------------------------------------------------------------------------
// Reconnector – context cancelled during backoff wait
// ---------------------------------------------------------------------------

func TestReconnector_TriggerReconnect_ContextCancelledDuringWait(t *testing.T) {
	// Nothing listening on this port so Connect always fails.
	ws := NewWSTransport(WSTransportConfig{
		URL:         "ws://127.0.0.1:1/nothing",
		DialTimeout: 50 * time.Millisecond,
	})

	r := NewReconnector(ws, ReconnectConfig{
		Enabled:     true,
		MaxAttempts: 0, // Unlimited.
		BaseWait:    500 * time.Millisecond,
		MaxWait:     500 * time.Millisecond,
	})

	// Cancel after a short time so the reconnector is in the backoff wait.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := r.TriggerReconnect(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("error = %v; want context.DeadlineExceeded", err)
	}
}

// ---------------------------------------------------------------------------
// Reconnector – backoff capped at MaxWait
// ---------------------------------------------------------------------------

func TestReconnector_BackoffCappedAtMaxWait(t *testing.T) {
	// Nothing listening -- Connect will always fail.
	ws := NewWSTransport(WSTransportConfig{
		URL:         "ws://127.0.0.1:1/nothing",
		DialTimeout: 10 * time.Millisecond,
	})

	r := NewReconnector(ws, ReconnectConfig{
		Enabled:     true,
		MaxAttempts: 5,
		BaseWait:    time.Millisecond,
		MaxWait:     5 * time.Millisecond, // Very small max to keep test fast.
	})

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.TriggerReconnect(ctx)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error")
	}
	// With 5 attempts and max 5ms backoff each + 10ms dial timeout each,
	// total should be well under 1 second.
	if elapsed > 2*time.Second {
		t.Errorf("reconnect took %v; expected under 2s with small backoff/MaxWait", elapsed)
	}
}

// ---------------------------------------------------------------------------
// Reconnector – Start and Stop lifecycle
// ---------------------------------------------------------------------------

func TestReconnector_StartStop_Lifecycle(t *testing.T) {
	ws := NewWSTransport(WSTransportConfig{URL: "ws://localhost:8080"})
	r := NewReconnector(ws, ReconnectConfig{Enabled: true})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start should set active to true.
	r.Start(ctx)

	r.mu.Lock()
	active := r.active
	r.mu.Unlock()
	if !active {
		t.Error("expected active to be true after Start")
	}

	// Stop should set active to false and clear cancelFunc.
	r.Stop()

	r.mu.Lock()
	active = r.active
	cf := r.cancelFunc
	r.mu.Unlock()
	if active {
		t.Error("expected active to be false after Stop")
	}
	if cf != nil {
		t.Error("expected cancelFunc to be nil after Stop")
	}

	// Start again should work.
	r.Start(ctx)

	r.mu.Lock()
	active = r.active
	r.mu.Unlock()
	if !active {
		t.Error("expected active to be true after re-Start")
	}
	r.Stop()
}

// ---------------------------------------------------------------------------
// WSTransport – Send ID generation is sequential
// ---------------------------------------------------------------------------

func TestWSTransport_SendIDGeneration(t *testing.T) {
	server, wsURL := wsTestServer(t, func(conn *gorilla.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	ws := NewWSTransport(WSTransportConfig{
		URL:          wsURL,
		PingInterval: 1 * time.Hour,
	})
	err := ws.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer ws.Close()

	id1, err := ws.Send(context.Background(), "m1", nil)
	if err != nil {
		t.Fatalf("Send 1 failed: %v", err)
	}
	id2, err := ws.Send(context.Background(), "m2", nil)
	if err != nil {
		t.Fatalf("Send 2 failed: %v", err)
	}
	id3, err := ws.Send(context.Background(), "m3", nil)
	if err != nil {
		t.Fatalf("Send 3 failed: %v", err)
	}

	if id2 != id1+1 || id3 != id2+1 {
		t.Errorf("expected sequential IDs, got %d, %d, %d", id1, id2, id3)
	}
}

// ---------------------------------------------------------------------------
// HTTPTransport – DoPublic with successful result parsing into nil result
// ---------------------------------------------------------------------------

func TestDoPublic_NilResultPtr(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`{"data":"ignored"}`)})
	}))
	defer server.Close()

	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       server.URL,
		RetryBaseWait: time.Millisecond,
	})
	// result is nil -- should not error.
	err := tr.DoPublic(context.Background(), "/ok", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// HTTPTransport – 4xx with non-JSON, non-API-error body
// ---------------------------------------------------------------------------

func TestDoPublic_4xxNonAPIErrorJSON(t *testing.T) {
	// JSON body but without an "error" field.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"message":"not found"}`))
	}))
	defer server.Close()

	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       server.URL,
		RetryBaseWait: time.Millisecond,
	})
	err := tr.DoPublic(context.Background(), "/missing", nil, nil)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "HTTP 404") {
		t.Errorf("error = %q; want to contain 'HTTP 404'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// HTTPTransport – API error in 200 response but nil error field
// ---------------------------------------------------------------------------

func TestDoPublic_200WithNilError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(apiResponse{
			Result: json.RawMessage(`42`),
			Error:  nil,
		})
	}))
	defer server.Close()

	tr := NewHTTPTransport(HTTPTransportConfig{
		BaseURL:       server.URL,
		RetryBaseWait: time.Millisecond,
	})
	var result int
	err := tr.DoPublic(context.Background(), "/ok", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 42 {
		t.Errorf("result = %d; want 42", result)
	}
}
