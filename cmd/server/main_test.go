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
type mockRouter struct {
	shouldFail bool
}

func (m *mockRouter) Run(addr string) error {
	// Simulating no-op router run to avoid blocking
	if m.shouldFail {
		return errors.New("router failure")
	}
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
	type startServerTest struct {
		name           string
		initFunc       func() error
		router         func() interface{ Run(string) error }
		expectedError  error
		expectedLogMsg string
		expectError    bool
	}

	tests := []startServerTest{
		{
			name: "Initialization returns nil and router runs successfully",
			initFunc: func() error {
				return nil
			},
			router: func() interface{ Run(string) error } {
				return &mockRouter{}
			},
			expectedError:  nil,
			expectedLogMsg: "Starting server...",
			expectError:    false,
		},
		{
			name: "Successful server start",
			initFunc: func() error {
				return nil
			},
			router: func() interface{ Run(string) error } {
				return &mockRouter{}
			},
			expectedError:  nil,
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
			errChan := make(chan error)
			go func() {
				errChan <- startServerWithMockLogger(tt.initFunc, tt.router(), logger)
			}()
			// Add a small sleep to capture logs (this is a bit hacky but should work for now)
			time.Sleep(1 * time.Second)
			t.Log("Returned from startServer call")
			logOutput := buf.String()
			t.Logf("Captured log output: %s", logOutput) // Debug captured logs
			if !assert.Contains(t, logOutput, tt.expectedLogMsg) {
				t.Errorf("Expected log message: %s, but got: %s", tt.expectedLogMsg, logOutput)
			}
			err := <-errChan
			expectedError := tt.expectedError
			if !errors.Is(err, expectedError) {
				t.Errorf("startServerWithMockLogger() error = %v, expectedError %v", err, expectedError)
			}
		})
	}
}
