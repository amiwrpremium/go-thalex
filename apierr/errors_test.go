package apierr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/amiwrpremium/go-thalex/apierr"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *apierr.APIError
		wantMsg string
	}{
		{
			name:    "BasicError",
			err:     &apierr.APIError{Code: 1001, Message: "invalid parameter"},
			wantMsg: "thalex: API error 1001: invalid parameter",
		},
		{
			name:    "ZeroCode",
			err:     &apierr.APIError{Code: 0, Message: "unknown error"},
			wantMsg: "thalex: API error 0: unknown error",
		},
		{
			name:    "EmptyMessage",
			err:     &apierr.APIError{Code: 500, Message: ""},
			wantMsg: "thalex: API error 500: ",
		},
		{
			name:    "NegativeCode",
			err:     &apierr.APIError{Code: -32600, Message: "invalid request"},
			wantMsg: "thalex: API error -32600: invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantMsg {
				t.Errorf("Error() = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestConnectionError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *apierr.ConnectionError
		wantMsg string
	}{
		{
			name:    "WithoutWrappedError",
			err:     &apierr.ConnectionError{Message: "connection refused"},
			wantMsg: "thalex: connection error: connection refused",
		},
		{
			name:    "WithWrappedError",
			err:     &apierr.ConnectionError{Message: "dial failed", Err: fmt.Errorf("timeout after 10s")},
			wantMsg: "thalex: connection error: dial failed: timeout after 10s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantMsg {
				t.Errorf("Error() = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestConnectionError_Unwrap(t *testing.T) {
	t.Run("WithWrappedError", func(t *testing.T) {
		inner := fmt.Errorf("inner error")
		connErr := &apierr.ConnectionError{Message: "msg", Err: inner}
		if connErr.Unwrap() != inner {
			t.Error("Unwrap() should return the inner error")
		}
	})

	t.Run("WithoutWrappedError", func(t *testing.T) {
		connErr := &apierr.ConnectionError{Message: "msg"}
		if connErr.Unwrap() != nil {
			t.Error("Unwrap() should return nil when no inner error")
		}
	})
}

func TestAuthError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *apierr.AuthError
		wantMsg string
	}{
		{
			name:    "WithoutWrappedError",
			err:     &apierr.AuthError{Message: "invalid key"},
			wantMsg: "thalex: auth error: invalid key",
		},
		{
			name:    "WithWrappedError",
			err:     &apierr.AuthError{Message: "token expired", Err: fmt.Errorf("JWT expired at 12:00")},
			wantMsg: "thalex: auth error: token expired: JWT expired at 12:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantMsg {
				t.Errorf("Error() = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestAuthError_Unwrap(t *testing.T) {
	t.Run("WithWrappedError", func(t *testing.T) {
		inner := fmt.Errorf("inner error")
		authErr := &apierr.AuthError{Message: "msg", Err: inner}
		if authErr.Unwrap() != inner {
			t.Error("Unwrap() should return the inner error")
		}
	})

	t.Run("WithoutWrappedError", func(t *testing.T) {
		authErr := &apierr.AuthError{Message: "msg"}
		if authErr.Unwrap() != nil {
			t.Error("Unwrap() should return nil when no inner error")
		}
	})
}

func TestTimeoutError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *apierr.TimeoutError
		wantMsg string
	}{
		{
			name:    "WithoutWrappedError",
			err:     &apierr.TimeoutError{Message: "request timed out"},
			wantMsg: "thalex: timeout: request timed out",
		},
		{
			name:    "WithWrappedError",
			err:     &apierr.TimeoutError{Message: "WS ping", Err: fmt.Errorf("context deadline exceeded")},
			wantMsg: "thalex: timeout: WS ping: context deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantMsg {
				t.Errorf("Error() = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestTimeoutError_Unwrap(t *testing.T) {
	t.Run("WithWrappedError", func(t *testing.T) {
		inner := fmt.Errorf("inner error")
		timeoutErr := &apierr.TimeoutError{Message: "msg", Err: inner}
		if timeoutErr.Unwrap() != inner {
			t.Error("Unwrap() should return the inner error")
		}
	})

	t.Run("WithoutWrappedError", func(t *testing.T) {
		timeoutErr := &apierr.TimeoutError{Message: "msg"}
		if timeoutErr.Unwrap() != nil {
			t.Error("Unwrap() should return nil when no inner error")
		}
	})
}

func TestIsAPIError(t *testing.T) {
	t.Run("WithAPIError", func(t *testing.T) {
		apiErr := &apierr.APIError{Code: 42, Message: "test"}
		got, ok := apierr.IsAPIError(apiErr)
		if !ok {
			t.Fatal("IsAPIError should return true for APIError")
		}
		if got != apiErr {
			t.Error("IsAPIError should return the same APIError pointer")
		}
		if got.Code != 42 {
			t.Errorf("Code = %d, want 42", got.Code)
		}
	})

	t.Run("WithWrappedAPIError", func(t *testing.T) {
		apiErr := &apierr.APIError{Code: 100, Message: "wrapped"}
		wrapped := fmt.Errorf("outer: %w", apiErr)
		got, ok := apierr.IsAPIError(wrapped)
		if !ok {
			t.Fatal("IsAPIError should return true for wrapped APIError")
		}
		if got.Code != 100 {
			t.Errorf("Code = %d, want 100", got.Code)
		}
	})

	t.Run("WithNonAPIError", func(t *testing.T) {
		plainErr := fmt.Errorf("plain error")
		got, ok := apierr.IsAPIError(plainErr)
		if ok {
			t.Fatal("IsAPIError should return false for non-APIError")
		}
		if got != nil {
			t.Error("IsAPIError should return nil for non-APIError")
		}
	})

	t.Run("WithConnectionError", func(t *testing.T) {
		connErr := &apierr.ConnectionError{Message: "conn"}
		_, ok := apierr.IsAPIError(connErr)
		if ok {
			t.Fatal("IsAPIError should return false for ConnectionError")
		}
	})

	t.Run("WithAuthError", func(t *testing.T) {
		authErr := &apierr.AuthError{Message: "auth"}
		_, ok := apierr.IsAPIError(authErr)
		if ok {
			t.Fatal("IsAPIError should return false for AuthError")
		}
	})

	t.Run("WithTimeoutError", func(t *testing.T) {
		timeoutErr := &apierr.TimeoutError{Message: "timeout"}
		_, ok := apierr.IsAPIError(timeoutErr)
		if ok {
			t.Fatal("IsAPIError should return false for TimeoutError")
		}
	})
}

func TestErrorsAs_APIError(t *testing.T) {
	apiErr := &apierr.APIError{Code: 1, Message: "test"}
	var target *apierr.APIError
	if !errors.As(apiErr, &target) {
		t.Error("errors.As should succeed for APIError")
	}
	if target.Code != 1 {
		t.Errorf("Code = %d, want 1", target.Code)
	}
}

func TestErrorsAs_ConnectionError(t *testing.T) {
	connErr := &apierr.ConnectionError{Message: "test", Err: fmt.Errorf("inner")}
	var target *apierr.ConnectionError
	if !errors.As(connErr, &target) {
		t.Error("errors.As should succeed for ConnectionError")
	}
	if target.Message != "test" {
		t.Errorf("Message = %q, want %q", target.Message, "test")
	}
}

func TestErrorsAs_AuthError(t *testing.T) {
	authErr := &apierr.AuthError{Message: "test", Err: fmt.Errorf("inner")}
	var target *apierr.AuthError
	if !errors.As(authErr, &target) {
		t.Error("errors.As should succeed for AuthError")
	}
	if target.Message != "test" {
		t.Errorf("Message = %q, want %q", target.Message, "test")
	}
}

func TestErrorsAs_TimeoutError(t *testing.T) {
	timeoutErr := &apierr.TimeoutError{Message: "test", Err: fmt.Errorf("inner")}
	var target *apierr.TimeoutError
	if !errors.As(timeoutErr, &target) {
		t.Error("errors.As should succeed for TimeoutError")
	}
	if target.Message != "test" {
		t.Errorf("Message = %q, want %q", target.Message, "test")
	}
}

func TestErrorsAs_WrappedConnectionError(t *testing.T) {
	inner := fmt.Errorf("root cause")
	connErr := &apierr.ConnectionError{Message: "conn fail", Err: inner}
	wrapped := fmt.Errorf("wrapping: %w", connErr)

	var target *apierr.ConnectionError
	if !errors.As(wrapped, &target) {
		t.Error("errors.As should find ConnectionError through wrapping")
	}
	if target.Message != "conn fail" {
		t.Errorf("Message = %q, want %q", target.Message, "conn fail")
	}
	if target.Unwrap() != inner {
		t.Error("Unwrap should return original inner error")
	}
}

func TestErrorsAs_WrappedAuthError(t *testing.T) {
	inner := fmt.Errorf("root cause")
	authErr := &apierr.AuthError{Message: "auth fail", Err: inner}
	wrapped := fmt.Errorf("wrapping: %w", authErr)

	var target *apierr.AuthError
	if !errors.As(wrapped, &target) {
		t.Error("errors.As should find AuthError through wrapping")
	}
}

func TestErrorsAs_WrappedTimeoutError(t *testing.T) {
	inner := fmt.Errorf("root cause")
	timeoutErr := &apierr.TimeoutError{Message: "timeout fail", Err: inner}
	wrapped := fmt.Errorf("wrapping: %w", timeoutErr)

	var target *apierr.TimeoutError
	if !errors.As(wrapped, &target) {
		t.Error("errors.As should find TimeoutError through wrapping")
	}
}

func TestAllErrorTypesImplementError(t *testing.T) {
	var _ error = (*apierr.APIError)(nil)
	var _ error = (*apierr.ConnectionError)(nil)
	var _ error = (*apierr.AuthError)(nil)
	var _ error = (*apierr.TimeoutError)(nil)
}
