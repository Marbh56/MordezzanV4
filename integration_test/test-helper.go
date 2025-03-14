package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mordezzanV4/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// createTestUser creates a test user for use in other tests
func createTestUser(t *testing.T, server *httptest.Server) int64 {
	userData := models.CreateUserInput{
		Username: "characteruser",
		Email:    "character@example.com",
		Password: "securepassword123",
	}
	payload, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("Failed to marshal user data: %v", err)
	}

	log.Println("Sending request to create test user...")
	req, err := http.NewRequest("POST", server.URL+"/users", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status %d, got %d. Response: %s", http.StatusCreated, resp.StatusCode, string(body))
	}

	var createdUser models.User
	if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	log.Printf("Created test user with ID: %d", createdUser.ID)
	return createdUser.ID
}

func createTestCharacter(t *testing.T, server *httptest.Server, userID int64) int64 {
	characterData := models.CreateCharacterInput{
		UserID:       userID,
		Name:         "Gandalf",
		Class:        "Wizard",
		Strength:     10,
		Dexterity:    12,
		Constitution: 14,
		Wisdom:       18,
		Intelligence: 16,
		Charisma:     14,
		HitPoints:    30,
	}

	payload, err := json.Marshal(characterData)
	if err != nil {
		t.Fatalf("Failed to marshal character data: %v", err)
	}

	req, err := http.NewRequest("POST", server.URL+"/characters", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var createdCharacter models.Character
	if err := json.NewDecoder(resp.Body).Decode(&createdCharacter); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return createdCharacter.ID
}

// NewTestLogger creates a new test logger for nicer test output
func NewTestLogger(t *testing.T) *TestLogger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := config.Build()
	sugar := logger.Sugar()
	return &TestLogger{
		t:      t,
		indent: 0,
		log:    sugar,
	}
}

// Color constants
var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

// TestLogger provides a structured logger for tests
type TestLogger struct {
	t      *testing.T
	indent int
	log    *zap.SugaredLogger
}

// Section logs a section header
func (l *TestLogger) Section(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	header := strings.Repeat("=", 80)
	l.t.Logf("\n%s%s%s\n%s %s %s\n%s%s%s\n",
		colorBold, header, colorReset,
		colorBold+colorBlue, msg, colorReset,
		colorBold, header, colorReset)
	l.log.Info("SECTION: " + msg)
}

// Step logs a test step
func (l *TestLogger) Step(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.t.Logf("%s➤ %s%s", colorBold+colorCyan, msg, colorReset)
	l.log.Infow("STEP", "message", msg)
}

// Info logs an info message
func (l *TestLogger) Info(format string, args ...interface{}) {
	indent := strings.Repeat("  ", l.indent)
	msg := fmt.Sprintf(format, args...)
	l.t.Logf("%s%s• %s%s", indent, colorBlue, msg, colorReset)
	l.log.Infow("INFO", "message", msg, "indent", l.indent)
}

// Success logs a success message
func (l *TestLogger) Success(format string, args ...interface{}) {
	indent := strings.Repeat("  ", l.indent)
	msg := fmt.Sprintf(format, args...)
	l.t.Logf("%s%s✓ %s%s", indent, colorGreen, msg, colorReset)
	l.log.Infow("SUCCESS", "message", msg, "indent", l.indent)
}

// Error logs an error message
func (l *TestLogger) Error(format string, args ...interface{}) {
	indent := strings.Repeat("  ", l.indent)
	msg := fmt.Sprintf(format, args...)
	l.t.Logf("%s%s✗ %s%s", indent, colorRed, msg, colorReset)
	l.log.Errorw("ERROR", "message", msg, "indent", l.indent)
}

// Warning logs a warning message
func (l *TestLogger) Warning(format string, args ...interface{}) {
	indent := strings.Repeat("  ", l.indent)
	msg := fmt.Sprintf(format, args...)
	l.t.Logf("%s%s! %s%s", indent, colorYellow, msg, colorReset)
	l.log.Warnw("WARNING", "message", msg, "indent", l.indent)
}

// Indent increases the indentation level
func (l *TestLogger) Indent() {
	l.indent++
}

// Outdent decreases the indentation level
func (l *TestLogger) Outdent() {
	if l.indent > 0 {
		l.indent--
	}
}

// CheckNoError checks if an error is nil and logs it
func (l *TestLogger) CheckNoError(err error, format string, args ...interface{}) bool {
	if err != nil {
		indent := strings.Repeat("  ", l.indent)
		msg := fmt.Sprintf(format, args...)
		l.t.Logf("%s%s✗ %s: %v%s", indent, colorRed, msg, err, colorReset)
		l.log.Errorw("ERROR", "message", msg, "error", err, "indent", l.indent)
		return false
	}
	return true
}

// TestWithRetry executes a test function with retries
func (l *TestLogger) TestWithRetry(name string, maxRetries int, testFn func() error) {
	l.Step("Running %s", name)
	l.Indent()
	var err error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			l.Warning("Retry attempt %d of %d", i, maxRetries)
		}
		err = testFn()
		if err == nil {
			l.Success("Test passed")
			l.Outdent()
			return
		}
	}
	l.Error("Test failed after %d attempts: %v", maxRetries, err)
	l.Outdent()
	l.t.Fatalf("Test %s failed after %d attempts: %v", name, maxRetries, err)
}

// Separator prints a separator line
func (l *TestLogger) Separator() {
	l.t.Logf("%s%s%s", colorCyan, strings.Repeat("-", 80), colorReset)
}

// ShouldUseColors determines if colored output should be used
func ShouldUseColors() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("GO_TEST_COLOR") == "0" {
		return false
	}
	fileInfo, _ := os.Stdout.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		return false
	}
	return true
}
