package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mordezzanV4/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TestUser struct {
	ID       int64
	Token    string
	Username string
	Email    string
}

func CreateTestUserWithAuth(t *testing.T, server *httptest.Server) *TestUser {
	randomSuffix := fmt.Sprintf("%d", time.Now().UnixNano()%100000)
	username := "testuser_" + randomSuffix
	email := "testuser_" + randomSuffix + "@example.com"

	userData := models.CreateUserInput{
		Username: username,
		Email:    email,
		Password: "securepassword123",
	}

	payload, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("Failed to marshal user data: %v", err)
	}

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

	// Generate JWT token for the user
	token := generateTestToken(t, createdUser.ID)

	return &TestUser{
		ID:       createdUser.ID,
		Token:    token,
		Username: username,
		Email:    email,
	}
}

func AuthenticatedRequest(t *testing.T, method, url string, body io.Reader, user *TestUser) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add authentication header
	req.Header.Set("Authorization", "Bearer "+user.Token)

	return req
}

// CreateTestCharacter creates a character for testing purposes
func CreateTestCharacter(t *testing.T, server *httptest.Server, user *TestUser) int64 {
	characterData := models.CreateCharacterInput{
		UserID:       user.ID,
		Name:         "Gandalf",
		Class:        "Wizard",
		Level:        5,
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

	req := AuthenticatedRequest(t, "POST", server.URL+"/characters", bytes.NewBuffer(payload), user)

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

// CreateTestWeapon creates a weapon for testing purposes
func CreateTestWeapon(t *testing.T, server *httptest.Server, user *TestUser) models.Weapon {
	weaponData := models.CreateWeaponInput{
		Name:        "Test Sword",
		Category:    "Melee",
		WeaponClass: 1,
		Cost:        15,
		Weight:      2,
		Damage:      "1d6",
		Properties:  "Versatile",
	}

	payload, err := json.Marshal(weaponData)
	if err != nil {
		t.Fatalf("Failed to marshal weapon data: %v", err)
	}

	var req *http.Request
	if user != nil {
		req = AuthenticatedRequest(t, "POST", server.URL+"/weapons", bytes.NewBuffer(payload), user)
	} else {
		req, err = http.NewRequest("POST", server.URL+"/weapons", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var createdWeapon models.Weapon
	if err := json.NewDecoder(resp.Body).Decode(&createdWeapon); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return createdWeapon
}

// CreateTestEquipment creates equipment for testing purposes
func CreateTestEquipment(t *testing.T, server *httptest.Server, user *TestUser) models.Equipment {
	equipmentData := models.CreateEquipmentInput{
		Name:        "Backpack",
		Description: "A sturdy leather backpack for adventuring",
		Cost:        2.0,
		Weight:      5,
	}

	payload, err := json.Marshal(equipmentData)
	if err != nil {
		t.Fatalf("Failed to marshal equipment data: %v", err)
	}

	var req *http.Request
	if user != nil {
		req = AuthenticatedRequest(t, "POST", server.URL+"/equipment", bytes.NewBuffer(payload), user)
	} else {
		req, err = http.NewRequest("POST", server.URL+"/equipment", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var createdEquipment models.Equipment
	if err := json.NewDecoder(resp.Body).Decode(&createdEquipment); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return createdEquipment
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

func generateTestToken(t *testing.T, userID int64) string {
	// Load the .env file to get the JWT_SECRET
	err := godotenv.Load()
	if err != nil {
		t.Logf("Warning: .env file not found, using app default secret")
	}

	// Get the JWT secret from environment or use the app's default
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "mordezzan_development_secret_key_not_for_production"
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    "user", // Add the role claim as expected by your auth middleware
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to generate JWT token: %v", err)
	}

	return tokenString
}

// Helper function to add auth token to requests
func addAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

func (l *TestLogger) AuthInfo(user *TestUser) {
	l.Info("Using authenticated user: ID=%d, Username=%s", user.ID, user.Username)
	l.Info("Authentication token is present: %v", user.Token != "")
}
