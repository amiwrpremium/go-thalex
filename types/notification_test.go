package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------- Notification JSON round-trip ----------

func TestNotification_JSONRoundTrip(t *testing.T) {
	n := types.Notification{
		ID:            "notif-123",
		Time:          1700000000.0,
		Category:      "trade",
		Title:         "Trade Executed",
		Message:       "Your buy order was filled",
		DisplayType:   enums.DisplayTypeSuccess,
		Read:          false,
		AccountName:   "main",
		AccountNumber: "ACC-001",
		Popup:         true,
	}

	data, err := json.Marshal(n)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Notification
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.ID != "notif-123" {
		t.Errorf("ID = %q, want %q", got.ID, "notif-123")
	}
	if got.Category != "trade" {
		t.Errorf("Category = %q, want %q", got.Category, "trade")
	}
	if got.Title != "Trade Executed" {
		t.Errorf("Title = %q, want %q", got.Title, "Trade Executed")
	}
	if got.DisplayType != enums.DisplayTypeSuccess {
		t.Errorf("DisplayType = %q, want %q", got.DisplayType, enums.DisplayTypeSuccess)
	}
	if got.Read != false {
		t.Errorf("Read = %v, want false", got.Read)
	}
	if got.Popup != true {
		t.Errorf("Popup = %v, want true", got.Popup)
	}
	if got.AccountName != "main" {
		t.Errorf("AccountName = %q, want %q", got.AccountName, "main")
	}
}

// ---------- NotificationsResult JSON round-trip ----------

func TestNotificationsResult_JSONRoundTrip(t *testing.T) {
	r := types.NotificationsResult{
		Notifications: []types.Notification{
			{ID: "n-1", Time: 1700000000.0, Category: "trade", Title: "T1", Message: "m1", DisplayType: enums.DisplayTypeInfo},
			{ID: "n-2", Time: 1700000001.0, Category: "system", Title: "T2", Message: "m2", DisplayType: enums.DisplayTypeWarning},
		},
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.NotificationsResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.Notifications) != 2 {
		t.Fatalf("len(Notifications) = %d, want 2", len(got.Notifications))
	}
	if got.Notifications[0].ID != "n-1" {
		t.Errorf("Notifications[0].ID = %q, want %q", got.Notifications[0].ID, "n-1")
	}
	if got.Notifications[1].DisplayType != enums.DisplayTypeWarning {
		t.Errorf("Notifications[1].DisplayType = %q, want %q", got.Notifications[1].DisplayType, enums.DisplayTypeWarning)
	}
}

func TestNotificationsResult_Empty(t *testing.T) {
	r := types.NotificationsResult{}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.NotificationsResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Notifications != nil {
		t.Errorf("Notifications = %v, want nil", got.Notifications)
	}
}

// ---------- NotificationPreferences JSON round-trip ----------

func TestNotificationPreferences_JSONRoundTrip(t *testing.T) {
	prefs := types.NotificationPreferences{
		"trade": types.NotificationPreference{
			Email: true,
			Inbox: true,
			Popup: true,
			SMS:   false,
		},
		"system": types.NotificationPreference{
			Email: false,
			Inbox: true,
			Popup: false,
			SMS:   false,
		},
	}

	data, err := json.Marshal(prefs)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.NotificationPreferences
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("len(prefs) = %d, want 2", len(got))
	}

	tradePref, ok := got["trade"]
	if !ok {
		t.Fatal("missing 'trade' key in preferences")
	}
	if !tradePref.Email {
		t.Error("trade.Email = false, want true")
	}
	if !tradePref.Inbox {
		t.Error("trade.Inbox = false, want true")
	}
	if !tradePref.Popup {
		t.Error("trade.Popup = false, want true")
	}

	sysPref := got["system"]
	if sysPref.Email {
		t.Error("system.Email = true, want false")
	}
}

// ---------- NotificationPreference JSON round-trip ----------

func TestNotificationPreference_JSONRoundTrip(t *testing.T) {
	p := types.NotificationPreference{
		Email: true,
		Inbox: true,
		Popup: false,
		SMS:   true,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.NotificationPreference
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Email != true {
		t.Errorf("Email = %v, want true", got.Email)
	}
	if got.Inbox != true {
		t.Errorf("Inbox = %v, want true", got.Inbox)
	}
	if got.Popup != false {
		t.Errorf("Popup = %v, want false", got.Popup)
	}
	if got.SMS != true {
		t.Errorf("SMS = %v, want true", got.SMS)
	}
}
