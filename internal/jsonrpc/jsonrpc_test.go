package jsonrpc_test

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
)

// ---------------------------------------------------------------------------
// Version constant
// ---------------------------------------------------------------------------

func TestVersion(t *testing.T) {
	if jsonrpc.Version != "2.0" {
		t.Fatalf("expected Version to be %q, got %q", "2.0", jsonrpc.Version)
	}
}

// ---------------------------------------------------------------------------
// NewRequest
// ---------------------------------------------------------------------------

func TestNewRequest(t *testing.T) {
	t.Run("sets JSONRPC, ID, method, and params", func(t *testing.T) {
		params := map[string]string{"key": "value"}
		req := jsonrpc.NewRequest(42, "test_method", params)

		if req.JSONRPC != "2.0" {
			t.Errorf("JSONRPC = %q; want %q", req.JSONRPC, "2.0")
		}
		if req.ID != 42 {
			t.Errorf("ID = %d; want %d", req.ID, 42)
		}
		if req.Method != "test_method" {
			t.Errorf("Method = %q; want %q", req.Method, "test_method")
		}
		if req.Params == nil {
			t.Fatal("Params should not be nil")
		}
	})

	t.Run("with nil params", func(t *testing.T) {
		req := jsonrpc.NewRequest(1, "no_params", nil)
		if req.Params != nil {
			t.Errorf("Params = %v; want nil", req.Params)
		}
	})

	t.Run("JSON marshaling produces correct output", func(t *testing.T) {
		params := map[string]int{"amount": 100}
		req := jsonrpc.NewRequest(7, "order/create", params)

		data, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}

		var raw map[string]json.RawMessage
		if err := json.Unmarshal(data, &raw); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}

		// Check jsonrpc field.
		var version string
		json.Unmarshal(raw["jsonrpc"], &version)
		if version != "2.0" {
			t.Errorf("jsonrpc = %q; want %q", version, "2.0")
		}

		// Check id field.
		var id uint64
		json.Unmarshal(raw["id"], &id)
		if id != 7 {
			t.Errorf("id = %d; want %d", id, 7)
		}

		// Check method field.
		var method string
		json.Unmarshal(raw["method"], &method)
		if method != "order/create" {
			t.Errorf("method = %q; want %q", method, "order/create")
		}

		// Check params field.
		var p map[string]int
		json.Unmarshal(raw["params"], &p)
		if p["amount"] != 100 {
			t.Errorf("params.amount = %d; want %d", p["amount"], 100)
		}
	})

	t.Run("JSON marshaling with nil params omits params field", func(t *testing.T) {
		req := jsonrpc.NewRequest(1, "ping", nil)
		data, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}

		var raw map[string]json.RawMessage
		json.Unmarshal(data, &raw)
		if _, exists := raw["params"]; exists {
			t.Error("expected params field to be omitted when nil")
		}
	})
}

// ---------------------------------------------------------------------------
// Response
// ---------------------------------------------------------------------------

func TestResponse_IsError(t *testing.T) {
	t.Run("returns true when Error is set", func(t *testing.T) {
		resp := jsonrpc.Response{
			Error: &jsonrpc.Error{Code: -32600, Message: "Invalid Request"},
		}
		if !resp.IsError() {
			t.Error("expected IsError() to return true")
		}
	})

	t.Run("returns false when Error is nil", func(t *testing.T) {
		resp := jsonrpc.Response{
			Result: json.RawMessage(`"ok"`),
		}
		if resp.IsError() {
			t.Error("expected IsError() to return false")
		}
	})
}

func TestResponse_Unmarshal(t *testing.T) {
	t.Run("unmarshal response with result", func(t *testing.T) {
		raw := `{"jsonrpc":"2.0","id":1,"result":{"price":42.5}}`
		var resp jsonrpc.Response
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.JSONRPC != "2.0" {
			t.Errorf("JSONRPC = %q; want %q", resp.JSONRPC, "2.0")
		}
		if resp.ID == nil || *resp.ID != 1 {
			t.Errorf("ID = %v; want 1", resp.ID)
		}
		if resp.Result == nil {
			t.Fatal("Result should not be nil")
		}
		if resp.Error != nil {
			t.Errorf("Error should be nil, got %v", resp.Error)
		}
	})

	t.Run("unmarshal response with error", func(t *testing.T) {
		raw := `{"jsonrpc":"2.0","id":2,"error":{"code":-32601,"message":"Method not found"}}`
		var resp jsonrpc.Response
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID == nil || *resp.ID != 2 {
			t.Errorf("ID = %v; want 2", resp.ID)
		}
		if resp.Error == nil {
			t.Fatal("expected Error to be set")
		}
		if resp.Error.Code != -32601 {
			t.Errorf("Error.Code = %d; want %d", resp.Error.Code, -32601)
		}
		if resp.Error.Message != "Method not found" {
			t.Errorf("Error.Message = %q; want %q", resp.Error.Message, "Method not found")
		}
	})

	t.Run("unmarshal response with null result", func(t *testing.T) {
		raw := `{"jsonrpc":"2.0","id":3,"result":null}`
		var resp jsonrpc.Response
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID == nil || *resp.ID != 3 {
			t.Errorf("ID = %v; want 3", resp.ID)
		}
		// json.RawMessage for "null" is the bytes `null`, but with omitempty
		// the Result may be nil or contain the literal null bytes.
		if resp.Error != nil {
			t.Error("Error should be nil")
		}
	})
}

// ---------------------------------------------------------------------------
// Error
// ---------------------------------------------------------------------------

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		message  string
		expected string
	}{
		{
			name:     "parse error",
			code:     -32700,
			message:  "Parse error",
			expected: "jsonrpc error -32700: Parse error",
		},
		{
			name:     "method not found",
			code:     -32601,
			message:  "Method not found",
			expected: "jsonrpc error -32601: Method not found",
		},
		{
			name:     "internal error",
			code:     -32603,
			message:  "Internal error",
			expected: "jsonrpc error -32603: Internal error",
		},
		{
			name:     "custom application error",
			code:     42,
			message:  "Insufficient funds",
			expected: "jsonrpc error 42: Insufficient funds",
		},
		{
			name:     "zero code",
			code:     0,
			message:  "Unknown",
			expected: "jsonrpc error 0: Unknown",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := &jsonrpc.Error{Code: tc.code, Message: tc.message}
			got := e.Error()
			if got != tc.expected {
				t.Errorf("Error() = %q; want %q", got, tc.expected)
			}
		})
	}
}

// Verify that *Error satisfies the error interface.
func TestError_ImplementsErrorInterface(t *testing.T) {
	var _ error = (*jsonrpc.Error)(nil)
}

// ---------------------------------------------------------------------------
// Notification
// ---------------------------------------------------------------------------

func TestNotification_Unmarshal(t *testing.T) {
	t.Run("unmarshal a notification JSON", func(t *testing.T) {
		raw := `{"jsonrpc":"2.0","method":"ticker","params":{"price":100}}`
		var notif jsonrpc.Notification
		if err := json.Unmarshal([]byte(raw), &notif); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if notif.JSONRPC != "2.0" {
			t.Errorf("JSONRPC = %q; want %q", notif.JSONRPC, "2.0")
		}
		if notif.Method != "ticker" {
			t.Errorf("Method = %q; want %q", notif.Method, "ticker")
		}
		if notif.Params == nil {
			t.Fatal("Params should not be nil")
		}
	})

	t.Run("notification has no ID field", func(t *testing.T) {
		raw := `{"jsonrpc":"2.0","method":"heartbeat"}`
		var notif jsonrpc.Notification
		if err := json.Unmarshal([]byte(raw), &notif); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Notification struct has no ID field; this is a structural assertion.
		// Marshal and verify no "id" key is emitted.
		data, err := json.Marshal(notif)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}
		var m map[string]json.RawMessage
		json.Unmarshal(data, &m)
		if _, exists := m["id"]; exists {
			t.Error("expected notification to have no id field")
		}
	})
}

// ---------------------------------------------------------------------------
// ParseMessage
// ---------------------------------------------------------------------------

func TestParseMessage(t *testing.T) {
	t.Run("parse a response with result", func(t *testing.T) {
		data := []byte(`{"jsonrpc":"2.0","id":10,"result":{"status":"ok"}}`)
		msg, err := jsonrpc.ParseMessage(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msg.Response == nil {
			t.Fatal("expected Response to be set")
		}
		if msg.Notification != nil {
			t.Error("expected Notification to be nil")
		}
		if msg.Response.ID == nil || *msg.Response.ID != 10 {
			t.Errorf("Response.ID = %v; want 10", msg.Response.ID)
		}
	})

	t.Run("parse a response with error", func(t *testing.T) {
		data := []byte(`{"jsonrpc":"2.0","id":11,"error":{"code":-1,"message":"fail"}}`)
		msg, err := jsonrpc.ParseMessage(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msg.Response == nil {
			t.Fatal("expected Response to be set")
		}
		if msg.Response.Error == nil {
			t.Fatal("expected Response.Error to be set")
		}
		if msg.Response.Error.Code != -1 {
			t.Errorf("Error.Code = %d; want %d", msg.Response.Error.Code, -1)
		}
	})

	t.Run("parse a notification", func(t *testing.T) {
		data := []byte(`{"jsonrpc":"2.0","method":"subscription","params":{"channel":"trades"}}`)
		msg, err := jsonrpc.ParseMessage(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msg.Notification == nil {
			t.Fatal("expected Notification to be set")
		}
		if msg.Response != nil {
			t.Error("expected Response to be nil")
		}
		if msg.Notification.Method != "subscription" {
			t.Errorf("Method = %q; want %q", msg.Notification.Method, "subscription")
		}
	})

	t.Run("parse invalid JSON returns error", func(t *testing.T) {
		data := []byte(`{not valid json}`)
		_, err := jsonrpc.ParseMessage(data)
		if err == nil {
			t.Fatal("expected an error for invalid JSON")
		}
	})

	t.Run("message with method AND id is treated as response", func(t *testing.T) {
		// If both method and ID are present, it falls through to response parsing.
		data := []byte(`{"jsonrpc":"2.0","id":5,"method":"something","result":"ok"}`)
		msg, err := jsonrpc.ParseMessage(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Because raw.ID is not nil, it's not a notification; it's a response.
		if msg.Response == nil {
			t.Fatal("expected Response to be set when id is present")
		}
		if msg.Notification != nil {
			t.Error("expected Notification to be nil when id is present")
		}
	})

	t.Run("message with no method and no id is treated as response", func(t *testing.T) {
		// Neither method nor ID; falls through to response branch.
		data := []byte(`{"jsonrpc":"2.0","result":"ok"}`)
		msg, err := jsonrpc.ParseMessage(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msg.Response == nil {
			t.Fatal("expected Response to be set")
		}
	})

	t.Run("empty JSON object parses as response", func(t *testing.T) {
		data := []byte(`{}`)
		msg, err := jsonrpc.ParseMessage(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msg.Response == nil {
			t.Fatal("expected Response to be set for empty object")
		}
	})
}

// ---------------------------------------------------------------------------
// IDGenerator
// ---------------------------------------------------------------------------

func TestIDGenerator_Next(t *testing.T) {
	t.Run("returns sequential IDs starting from 1", func(t *testing.T) {
		var gen jsonrpc.IDGenerator
		for i := uint64(1); i <= 10; i++ {
			got := gen.Next()
			if got != i {
				t.Errorf("Next() = %d; want %d", got, i)
			}
		}
	})

	t.Run("concurrent safety - no duplicates", func(t *testing.T) {
		var gen jsonrpc.IDGenerator
		const goroutines = 100
		const idsPerGoroutine = 100

		results := make(chan uint64, goroutines*idsPerGoroutine)
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for g := 0; g < goroutines; g++ {
			go func() {
				defer wg.Done()
				for i := 0; i < idsPerGoroutine; i++ {
					results <- gen.Next()
				}
			}()
		}

		wg.Wait()
		close(results)

		seen := make(map[uint64]bool, goroutines*idsPerGoroutine)
		for id := range results {
			if seen[id] {
				t.Fatalf("duplicate ID: %d", id)
			}
			seen[id] = true
		}

		if len(seen) != goroutines*idsPerGoroutine {
			t.Errorf("expected %d unique IDs, got %d", goroutines*idsPerGoroutine, len(seen))
		}
	})

	t.Run("independent generators have independent counters", func(t *testing.T) {
		var gen1, gen2 jsonrpc.IDGenerator
		id1 := gen1.Next()
		id2 := gen2.Next()
		if id1 != 1 {
			t.Errorf("gen1.Next() = %d; want 1", id1)
		}
		if id2 != 1 {
			t.Errorf("gen2.Next() = %d; want 1", id2)
		}
	})
}
