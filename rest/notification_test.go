package rest

import (
	"context"
	"net/http"
	"testing"

	"github.com/amiwrpremium/go-thalex/types"
)

func TestNotificationsInbox_NoLimit(t *testing.T) {
	expected := types.NotificationsResult{
		Notifications: []types.Notification{
			{ID: "notif-001", Title: "Test", Message: "Hello"},
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/private/notifications_inbox" {
			t.Errorf("expected path /private/notifications_inbox, got %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query params, got %s", r.URL.RawQuery)
		}
		w.Write(wrapResult(t, expected))
	})

	result, err := c.NotificationsInbox(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Notifications) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(result.Notifications))
	}
	if result.Notifications[0].ID != "notif-001" {
		t.Errorf("expected ID=notif-001, got %s", result.Notifications[0].ID)
	}
}

func TestNotificationsInbox_WithLimit(t *testing.T) {
	expected := types.NotificationsResult{
		Notifications: []types.Notification{
			{ID: "notif-001"},
			{ID: "notif-002"},
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("limit") != "5" {
			t.Errorf("expected limit=5, got %s", q.Get("limit"))
		}
		w.Write(wrapResult(t, expected))
	})

	limit := 5
	result, err := c.NotificationsInbox(context.Background(), &limit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Notifications) != 2 {
		t.Fatalf("expected 2 notifications, got %d", len(result.Notifications))
	}
}

func TestMarkNotificationAsRead_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/private/mark_inbox_notification_as_read" {
			t.Errorf("expected path /private/mark_inbox_notification_as_read, got %s", r.URL.Path)
		}
		w.Write([]byte(`{"result":null}`))
	})

	err := c.MarkNotificationAsRead(context.Background(), "notif-001", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
