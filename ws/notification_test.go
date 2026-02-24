package ws

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
)

// ---------------------------------------------------------------------------
// NotificationsInbox
// ---------------------------------------------------------------------------

func TestNotificationsInbox_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/notifications_inbox" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"notifications":[{"id":"n1","time":1.0,"category":"trade","title":"Trade Executed","message":"Your trade was executed","display_type":"popup","read":false,"popup":true}]}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.NotificationsInbox(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Notifications) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(result.Notifications))
	}
	if result.Notifications[0].ID != "n1" {
		t.Errorf("expected id=n1, got %q", result.Notifications[0].ID)
	}
}

func TestNotificationsInbox_WithLimit(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`{"notifications":[]}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limit := 5
	result, err := c.NotificationsInbox(ctx, &limit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Notifications == nil {
		t.Error("expected non-nil notifications slice")
	}
}

// ---------------------------------------------------------------------------
// MarkNotificationAsRead
// ---------------------------------------------------------------------------

func TestMarkNotificationAsRead_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/mark_inbox_notification_as_read" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.MarkNotificationAsRead(ctx, "n1", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarkNotificationAsRead_Unread(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.MarkNotificationAsRead(ctx, "n1", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
