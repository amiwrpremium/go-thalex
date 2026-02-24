package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------- SystemInfo JSON round-trip ----------

func TestSystemInfo_JSONRoundTrip(t *testing.T) {
	bannerID := 1
	si := types.SystemInfo{
		Environment: "production",
		APIVersion:  "1.0.0",
		Banners: []types.Banner{
			{
				ID:       &bannerID,
				Time:     1700000000.0,
				Severity: enums.SeverityWarning,
				Title:    "Maintenance",
				Message:  "Scheduled maintenance at midnight",
			},
		},
	}

	data, err := json.Marshal(si)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.SystemInfo
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Environment != "production" {
		t.Errorf("Environment = %q, want %q", got.Environment, "production")
	}
	if got.APIVersion != "1.0.0" {
		t.Errorf("APIVersion = %q, want %q", got.APIVersion, "1.0.0")
	}
	if len(got.Banners) != 1 {
		t.Fatalf("len(Banners) = %d, want 1", len(got.Banners))
	}
	if got.Banners[0].ID == nil || *got.Banners[0].ID != 1 {
		t.Errorf("Banners[0].ID = %v, want 1", got.Banners[0].ID)
	}
	if got.Banners[0].Severity != enums.SeverityWarning {
		t.Errorf("Banners[0].Severity = %q, want %q", got.Banners[0].Severity, enums.SeverityWarning)
	}
	if got.Banners[0].Title != "Maintenance" {
		t.Errorf("Banners[0].Title = %q, want %q", got.Banners[0].Title, "Maintenance")
	}
}

func TestSystemInfo_EmptyBanners(t *testing.T) {
	si := types.SystemInfo{
		Environment: "testnet",
	}

	data, err := json.Marshal(si)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.SystemInfo
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Environment != "testnet" {
		t.Errorf("Environment = %q, want %q", got.Environment, "testnet")
	}
}

// ---------- Banner JSON round-trip ----------

func TestBanner_JSONRoundTrip(t *testing.T) {
	id := 42
	b := types.Banner{
		ID:       &id,
		Time:     1700000000.0,
		Severity: enums.SeverityCritical,
		Title:    "Outage",
		Message:  "Trading is halted",
	}

	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Banner
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.ID == nil || *got.ID != 42 {
		t.Errorf("ID = %v, want 42", got.ID)
	}
	if got.Severity != enums.SeverityCritical {
		t.Errorf("Severity = %q, want %q", got.Severity, enums.SeverityCritical)
	}
	if got.Message != "Trading is halted" {
		t.Errorf("Message = %q, want %q", got.Message, "Trading is halted")
	}
}

func TestBanner_NilID(t *testing.T) {
	b := types.Banner{
		Time:     1700000000.0,
		Severity: enums.SeverityInfo,
		Message:  "Info banner",
	}

	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	if _, ok := raw["id"]; ok {
		t.Error("expected id to be omitted when nil")
	}
}

// ---------- SystemEvent JSON round-trip ----------

func TestSystemEvent_JSONRoundTrip(t *testing.T) {
	se := types.SystemEvent{
		Event: enums.SystemEventTypeReconnect,
	}

	data, err := json.Marshal(se)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.SystemEvent
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Event != enums.SystemEventTypeReconnect {
		t.Errorf("Event = %q, want %q", got.Event, enums.SystemEventTypeReconnect)
	}
}
