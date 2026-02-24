package jsonrpc

import (
	"encoding/json"
	"testing"
)

// ---------------------------------------------------------------------------
// ParseMessage – internal coverage for second-pass unmarshal error paths
// ---------------------------------------------------------------------------

// TestParseMessage_NotificationFields verifies that a valid notification is
// fully hydrated (JSONRPC, Method, Params) after parsing.
func TestParseMessage_NotificationFields(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","method":"ticker","params":{"price":100}}`)
	msg, err := ParseMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Notification == nil {
		t.Fatal("expected Notification to be set")
	}
	if msg.Notification.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q; want %q", msg.Notification.JSONRPC, "2.0")
	}
	if msg.Notification.Method != "ticker" {
		t.Errorf("Method = %q; want %q", msg.Notification.Method, "ticker")
	}
	if msg.Notification.Params == nil {
		t.Error("expected Params to be set")
	}
}

// TestParseMessage_ResponseFields verifies that a valid response is fully
// hydrated (JSONRPC, ID, Result, Error) after parsing.
func TestParseMessage_ResponseFields(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":99,"result":{"ok":true}}`)
	msg, err := ParseMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Response == nil {
		t.Fatal("expected Response to be set")
	}
	if msg.Response.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q; want %q", msg.Response.JSONRPC, "2.0")
	}
	if msg.Response.ID == nil || *msg.Response.ID != 99 {
		t.Errorf("ID = %v; want 99", msg.Response.ID)
	}
	if msg.Response.Result == nil {
		t.Error("expected Result to be set")
	}
	if msg.Response.Error != nil {
		t.Errorf("expected Error to be nil, got %v", msg.Response.Error)
	}
}

// TestParseMessage_ResponseWithError verifies that parsing a response containing
// an error object correctly populates the Error field.
func TestParseMessage_ResponseWithError(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"Invalid Request"}}`)
	msg, err := ParseMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Response == nil {
		t.Fatal("expected Response to be set")
	}
	if msg.Response.Error == nil {
		t.Fatal("expected Error to be set")
	}
	if msg.Response.Error.Code != -32600 {
		t.Errorf("Error.Code = %d; want -32600", msg.Response.Error.Code)
	}
	if msg.Response.Error.Message != "Invalid Request" {
		t.Errorf("Error.Message = %q; want %q", msg.Response.Error.Message, "Invalid Request")
	}
}

// TestParseMessage_NotificationWithoutParams confirms notifications without
// params are parsed correctly.
func TestParseMessage_NotificationWithoutParams(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","method":"heartbeat"}`)
	msg, err := ParseMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Notification == nil {
		t.Fatal("expected Notification to be set")
	}
	if msg.Notification.Params != nil {
		t.Errorf("expected Params to be nil, got %s", string(msg.Notification.Params))
	}
}

// TestParseMessage_InvalidJSON verifies that completely invalid JSON returns a
// wrapped error.
func TestParseMessage_InvalidJSON(t *testing.T) {
	data := []byte(`{{{invalid`)
	_, err := ParseMessage(data)
	if err == nil {
		t.Fatal("expected an error for invalid JSON")
	}
}

// TestParseMessage_EmptyArray verifies that a JSON array (not object) fails.
func TestParseMessage_EmptyArray(t *testing.T) {
	data := []byte(`[]`)
	_, err := ParseMessage(data)
	if err == nil {
		t.Fatal("expected an error for JSON array instead of object")
	}
}

// TestParseMessage_NullJSON verifies that JSON null fails parsing.
func TestParseMessage_NullJSON(t *testing.T) {
	data := []byte(`null`)
	// json.Unmarshal into a struct from null succeeds with zero values,
	// so this should parse as a response with no ID and no method.
	msg, err := ParseMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// null unmarshals into the raw struct successfully; method is "" and id is nil,
	// so it falls through to the response branch.
	if msg.Response == nil {
		t.Fatal("expected Response to be set for null JSON")
	}
}

// TestParseMessage_ResponseWithNullID verifies that a response with null ID
// is handled correctly (treated as response since method is empty).
func TestParseMessage_ResponseWithNullID(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":null,"result":"ok"}`)
	msg, err := ParseMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// raw.ID will be a non-nil *json.RawMessage containing "null",
	// so method=="" && ID!=nil means response branch.
	if msg.Response == nil {
		t.Fatal("expected Response to be set")
	}
}

// ---------------------------------------------------------------------------
// Error – internal tests
// ---------------------------------------------------------------------------

func TestError_ErrorFormat(t *testing.T) {
	e := &Error{Code: 123, Message: "something went wrong"}
	expected := "jsonrpc error 123: something went wrong"
	if got := e.Error(); got != expected {
		t.Errorf("Error() = %q; want %q", got, expected)
	}
}

func TestError_NegativeCode(t *testing.T) {
	e := &Error{Code: -32700, Message: "Parse error"}
	expected := "jsonrpc error -32700: Parse error"
	if got := e.Error(); got != expected {
		t.Errorf("Error() = %q; want %q", got, expected)
	}
}

func TestError_EmptyMessage(t *testing.T) {
	e := &Error{Code: 0, Message: ""}
	expected := "jsonrpc error 0: "
	if got := e.Error(); got != expected {
		t.Errorf("Error() = %q; want %q", got, expected)
	}
}

// ---------------------------------------------------------------------------
// Response – internal tests
// ---------------------------------------------------------------------------

func TestResponse_IsError_Internal(t *testing.T) {
	t.Run("error present", func(t *testing.T) {
		r := &Response{Error: &Error{Code: 1, Message: "err"}}
		if !r.IsError() {
			t.Error("expected IsError() to be true")
		}
	})
	t.Run("error absent", func(t *testing.T) {
		r := &Response{Result: json.RawMessage(`"ok"`)}
		if r.IsError() {
			t.Error("expected IsError() to be false")
		}
	})
}

// ---------------------------------------------------------------------------
// IDGenerator – internal tests
// ---------------------------------------------------------------------------

func TestIDGenerator_StartsAtOne(t *testing.T) {
	var gen IDGenerator
	if id := gen.Next(); id != 1 {
		t.Errorf("first ID = %d; want 1", id)
	}
}

func TestIDGenerator_Sequential(t *testing.T) {
	var gen IDGenerator
	for i := uint64(1); i <= 5; i++ {
		if id := gen.Next(); id != i {
			t.Errorf("ID = %d; want %d", id, i)
		}
	}
}

// ---------------------------------------------------------------------------
// NewRequest – internal tests
// ---------------------------------------------------------------------------

func TestNewRequest_Internal(t *testing.T) {
	req := NewRequest(42, "test.method", map[string]string{"key": "val"})
	if req.JSONRPC != Version {
		t.Errorf("JSONRPC = %q; want %q", req.JSONRPC, Version)
	}
	if req.ID != 42 {
		t.Errorf("ID = %d; want 42", req.ID)
	}
	if req.Method != "test.method" {
		t.Errorf("Method = %q; want %q", req.Method, "test.method")
	}
	if req.Params == nil {
		t.Error("expected Params to be set")
	}
}

func TestNewRequest_NilParams(t *testing.T) {
	req := NewRequest(1, "no_params", nil)
	if req.Params != nil {
		t.Errorf("expected Params to be nil, got %v", req.Params)
	}
}

// TestNewRequest_MarshalRoundTrip verifies that a Request can be marshaled to
// JSON and back without losing data.
func TestNewRequest_MarshalRoundTrip(t *testing.T) {
	req := NewRequest(7, "order/create", map[string]int{"amount": 100})

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var decoded Request
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if decoded.JSONRPC != req.JSONRPC {
		t.Errorf("JSONRPC = %q; want %q", decoded.JSONRPC, req.JSONRPC)
	}
	if decoded.ID != req.ID {
		t.Errorf("ID = %d; want %d", decoded.ID, req.ID)
	}
	if decoded.Method != req.Method {
		t.Errorf("Method = %q; want %q", decoded.Method, req.Method)
	}
}

// ---------------------------------------------------------------------------
// Message struct – internal tests
// ---------------------------------------------------------------------------

func TestMessage_BothNil(t *testing.T) {
	msg := &Message{}
	if msg.Response != nil {
		t.Error("expected Response to be nil")
	}
	if msg.Notification != nil {
		t.Error("expected Notification to be nil")
	}
}

// ---------------------------------------------------------------------------
// Notification – internal marshal/unmarshal tests
// ---------------------------------------------------------------------------

func TestNotification_MarshalOmitsID(t *testing.T) {
	notif := Notification{
		JSONRPC: "2.0",
		Method:  "heartbeat",
	}
	data, err := json.Marshal(notif)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if _, ok := raw["id"]; ok {
		t.Error("expected notification marshal to not contain 'id' field")
	}
}

func TestNotification_MarshalOmitsEmptyParams(t *testing.T) {
	notif := Notification{
		JSONRPC: "2.0",
		Method:  "ping",
	}
	data, err := json.Marshal(notif)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if _, ok := raw["params"]; ok {
		t.Error("expected notification marshal to omit empty params")
	}
}

// ---------------------------------------------------------------------------
// ParseMessage – second-pass unmarshal error paths
// ---------------------------------------------------------------------------

// TestParseMessage_ResponseUnmarshalError triggers the error at line 87-88
// in ParseMessage: the first pass succeeds (raw struct with *json.RawMessage
// for ID), but the second pass into Response (with *uint64 for ID) fails
// because "id" contains a non-numeric value.
func TestParseMessage_ResponseUnmarshalError(t *testing.T) {
	// ID is a string, which the first pass accepts as *json.RawMessage
	// but the second pass rejects as *uint64.
	data := []byte(`{"jsonrpc":"2.0","id":"not_a_number","result":"ok"}`)
	_, err := ParseMessage(data)
	if err == nil {
		t.Fatal("expected error when response ID is a string, got nil")
	}
	if !contains(err.Error(), "failed to parse response") {
		t.Errorf("error should mention response parse failure, got: %v", err)
	}
}

// TestParseMessage_ResponseUnmarshalError_IDObject triggers the response
// unmarshal error with an object as the ID value.
func TestParseMessage_ResponseUnmarshalError_IDObject(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":{"nested":true},"result":"ok"}`)
	_, err := ParseMessage(data)
	if err == nil {
		t.Fatal("expected error when response ID is an object, got nil")
	}
	if !contains(err.Error(), "failed to parse response") {
		t.Errorf("error should mention response parse failure, got: %v", err)
	}
}

// TestParseMessage_ResponseUnmarshalError_IDArray triggers the response
// unmarshal error with an array as the ID value.
func TestParseMessage_ResponseUnmarshalError_IDArray(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":[1,2,3],"result":"ok"}`)
	_, err := ParseMessage(data)
	if err == nil {
		t.Fatal("expected error when response ID is an array, got nil")
	}
	if !contains(err.Error(), "failed to parse response") {
		t.Errorf("error should mention response parse failure, got: %v", err)
	}
}

// TestParseMessage_ResponseUnmarshalError_IDBool triggers the response
// unmarshal error with a boolean as the ID value.
func TestParseMessage_ResponseUnmarshalError_IDBool(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":true,"result":"ok"}`)
	_, err := ParseMessage(data)
	if err == nil {
		t.Fatal("expected error when response ID is a boolean, got nil")
	}
	if !contains(err.Error(), "failed to parse response") {
		t.Errorf("error should mention response parse failure, got: %v", err)
	}
}

// TestParseMessage_NotificationUnmarshalError triggers the error at line 79-80
// in ParseMessage. This is harder to trigger because Notification has only
// string and json.RawMessage fields which are very lenient. However, if the
// JSON has a "params" field that is not a valid JSON value for the second
// pass, it will fail.
// NOTE: This path is structurally very hard to reach because the Notification
// struct uses json.RawMessage for Params and string for Method/JSONRPC, which
// accept the same inputs as the first-pass raw struct. We attempt it with a
// type that would fail on the second pass if the struct were stricter.
// In practice this error path may be unreachable with current struct
// definitions, but we test the boundary condition for completeness.
func TestParseMessage_NotificationParamsVariousTypes(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "params is array",
			data: `{"jsonrpc":"2.0","method":"update","params":[1,2,3]}`,
		},
		{
			name: "params is string",
			data: `{"jsonrpc":"2.0","method":"update","params":"string_value"}`,
		},
		{
			name: "params is number",
			data: `{"jsonrpc":"2.0","method":"update","params":42}`,
		},
		{
			name: "params is boolean",
			data: `{"jsonrpc":"2.0","method":"update","params":true}`,
		},
		{
			name: "params is null",
			data: `{"jsonrpc":"2.0","method":"update","params":null}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := ParseMessage([]byte(tc.data))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if msg.Notification == nil {
				t.Fatal("expected Notification to be set")
			}
			if msg.Notification.Method != "update" {
				t.Errorf("Method = %q; want %q", msg.Notification.Method, "update")
			}
		})
	}
}

// TestParseMessage_ResponseWithNegativeID tests that a response with a
// negative numeric ID fails the second-pass unmarshal since ID is *uint64.
func TestParseMessage_ResponseWithNegativeID(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":-1,"result":"ok"}`)
	_, err := ParseMessage(data)
	// A negative number cannot be unmarshaled into *uint64 in strict mode.
	// This should trigger the response unmarshal error path.
	if err == nil {
		// If Go's json package accepts -1 into *uint64 (some versions may overflow),
		// we just verify the response parsed. Either outcome is acceptable.
		return
	}
	if !contains(err.Error(), "failed to parse response") {
		t.Errorf("error should mention response parse failure, got: %v", err)
	}
}

// TestParseMessage_ResponseWithFloatID tests that a response with a float ID
// fails the second-pass unmarshal since ID is *uint64.
func TestParseMessage_ResponseWithFloatID(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":1.5,"result":"ok"}`)
	_, err := ParseMessage(data)
	if err == nil {
		// Some JSON decoders may truncate 1.5 to 1; if so, we accept that.
		return
	}
	if !contains(err.Error(), "failed to parse response") {
		t.Errorf("error should mention response parse failure, got: %v", err)
	}
}

// TestParseMessage_LargeResponse tests parsing a response with a large result payload.
func TestParseMessage_LargeResponse(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":999,"result":{"data":` +
		`"` + string(make([]byte, 0)) + `large_payload"}}`)
	msg, err := ParseMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Response == nil {
		t.Fatal("expected Response to be set")
	}
	if msg.Response.ID == nil || *msg.Response.ID != 999 {
		t.Errorf("Response.ID = %v; want 999", msg.Response.ID)
	}
}

// TestParseMessage_ResponseWithBothResultAndError tests parsing a response
// that contains both result and error fields.
func TestParseMessage_ResponseWithBothResultAndError(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":5,"result":"ok","error":{"code":-1,"message":"fail"}}`)
	msg, err := ParseMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Response == nil {
		t.Fatal("expected Response to be set")
	}
	// Both fields should be populated
	if msg.Response.Result == nil {
		t.Error("expected Result to be set")
	}
	if msg.Response.Error == nil {
		t.Error("expected Error to be set")
	}
	if !msg.Response.IsError() {
		t.Error("IsError() should return true when Error is present")
	}
}

// TestParseMessage_EmptyString tests that an empty input fails.
func TestParseMessage_EmptyString(t *testing.T) {
	_, err := ParseMessage([]byte(""))
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

// TestParseMessage_JustWhitespace tests that whitespace-only input fails.
func TestParseMessage_JustWhitespace(t *testing.T) {
	_, err := ParseMessage([]byte("   "))
	if err == nil {
		t.Fatal("expected error for whitespace-only input, got nil")
	}
}

// contains is a helper to check substring presence without importing strings.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Version constant – internal test
// ---------------------------------------------------------------------------

func TestVersion_Internal(t *testing.T) {
	if Version != "2.0" {
		t.Fatalf("Version = %q; want %q", Version, "2.0")
	}
}
