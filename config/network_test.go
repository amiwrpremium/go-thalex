package config_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/amiwrpremium/go-thalex/config"
)

func TestNetworkBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		network config.Network
		want    string
	}{
		{
			name:    "Production",
			network: config.Production,
			want:    "https://thalex.com/api/v2",
		},
		{
			name:    "Testnet",
			network: config.Testnet,
			want:    "https://testnet.thalex.com/api/v2",
		},
		{
			name:    "InvalidNetworkFallsBackToProduction",
			network: config.Network(99),
			want:    "https://thalex.com/api/v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.network.BaseURL()
			if got != tt.want {
				t.Errorf("BaseURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNetworkWebSocketURL(t *testing.T) {
	tests := []struct {
		name    string
		network config.Network
		want    string
	}{
		{
			name:    "Production",
			network: config.Production,
			want:    "wss://thalex.com/ws/api/v2",
		},
		{
			name:    "Testnet",
			network: config.Testnet,
			want:    "wss://testnet.thalex.com/ws/api/v2",
		},
		{
			name:    "InvalidNetworkFallsBackToProduction",
			network: config.Network(99),
			want:    "wss://thalex.com/ws/api/v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.network.WebSocketURL()
			if got != tt.want {
				t.Errorf("WebSocketURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNetworkString(t *testing.T) {
	tests := []struct {
		name    string
		network config.Network
		want    string
	}{
		{
			name:    "Production",
			network: config.Production,
			want:    "production",
		},
		{
			name:    "Testnet",
			network: config.Testnet,
			want:    "testnet",
		},
		{
			name:    "InvalidNetworkFallsBackToProduction",
			network: config.Network(42),
			want:    "production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.network.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVersionConstant(t *testing.T) {
	if config.Version == "" {
		t.Fatal("Version constant should not be empty")
	}

	// Validate semver format: MAJOR.MINOR.PATCH or MAJOR.MINOR.PATCH-prerelease
	core := config.Version
	if idx := strings.IndexByte(core, '-'); idx != -1 {
		core = core[:idx]
	}
	segments := strings.Split(core, ".")
	if len(segments) != 3 {
		t.Fatalf("Version %q does not have 3 dot-separated segments", config.Version)
	}
	for i, seg := range segments {
		if _, err := strconv.Atoi(seg); err != nil {
			t.Errorf("Version segment %d (%q) is not a valid integer", i, seg)
		}
	}
}

func TestUserAgentConstant(t *testing.T) {
	want := "go-thalex/" + config.Version
	if config.UserAgent != want {
		t.Errorf("UserAgent = %q, want %q", config.UserAgent, want)
	}
}

func TestProductionIsZeroValue(t *testing.T) {
	var n config.Network
	if n != config.Production {
		t.Errorf("zero-value Network should be Production, got %v", n)
	}
}
