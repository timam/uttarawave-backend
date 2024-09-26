package server

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

// mockRouter to prevent blocking in tests
type mockRouter struct{}

func (m *mockRouter) Run(addr string) error {
	// Simulating no-op router run to avoid blocking
	return nil
}

// Test-specific wrapper to simulate log.Fatalf behavior (without process exit)
func mockLogFatalf(logger *log.Logger, format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func startServerWithMockLogger(initFunc func() error, router interface{ Run(string) error }, logger *log.Logger) error {
	err := initFunc()
	if err != nil {
		mockLogFatalf(logger, "Initialization failed: %v", err)
		return err
	}

	logger.Println("Starting server...")
	err = router.Run(":8080")
	if err != nil {
		mockLogFatalf(logger, "Failed to start server: %v", err)
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestStartServer(t *testing.T) {
	tests := []struct {
		name           string
		initFunc       func() error
		expectedLogMsg string
		expectError    bool
	}{
		{
			name: "Initialization failure",
			initFunc: func() error {
				return errors.New("initialization failed")
			},
			expectedLogMsg: "Initialization failed: initialization failed",
			expectError:    true,
		},
		{
			name: "Successful server start",
			initFunc: func() error {
				return nil
			},
			expectedLogMsg: "Starting server...",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Starting test:", tt.name) // Initial debug statement

			// Capture log output
			var buf bytes.Buffer
			logger := log.New(&buf, "", log.LstdFlags)
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			t.Log("Calling startServer with injected initFunc")
			go func() {
				_ = startServerWithMockLogger(tt.initFunc, &mockRouter{}, logger)
			}()

			// Add a small sleep to capture logs (this is a bit hacky but should work for now)
			time.Sleep(1 * time.Second)

			t.Log("Returned from startServer call")

			logOutput := buf.String()
			t.Logf("Captured log output: %s", logOutput) // Debug captured logs
			if !assert.Contains(t, logOutput, tt.expectedLogMsg) {
				t.Errorf("Expected log message: %s, but got: %s", tt.expectedLogMsg, logOutput)
			}
		})
	}
}
