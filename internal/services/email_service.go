package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
)

// EmailConfig holds the configuration for sending emails
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	FromAddr string
	// Set to true if you want to skip TLS verification (not recommended for production)
	InsecureSkipVerify bool
}

// EmailService handles email sending operations
type EmailService struct {
	config    EmailConfig
	templates *template.Template
}

// NewEmailService creates a new EmailService with the given configuration
func NewEmailService(config EmailConfig, templatesDir string) (*EmailService, error) {
	// Load email templates from the templates directory
	templates, err := template.ParseGlob(filepath.Join(templatesDir, "email/*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse email templates: %w", err)
	}

	return &EmailService{
		config:    config,
		templates: templates,
	}, nil
}

// NewEmailServiceFromEnv creates a new EmailService using environment variables
func NewEmailServiceFromEnv(templatesDir string) (*EmailService, error) {
	config := EmailConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     587, // Default port for TLS
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		FromName: os.Getenv("EMAIL_FROM_NAME"),
		FromAddr: os.Getenv("EMAIL_FROM_ADDR"),
	}

	// Parse port from environment if provided
	if portStr := os.Getenv("EMAIL_PORT"); portStr != "" {
		var port int
		if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
			return nil, fmt.Errorf("invalid EMAIL_PORT: %w", err)
		}
		config.Port = port
	}

	// Check for required configuration
	if config.Host == "" || config.Username == "" || config.Password == "" || config.FromAddr == "" {
		return nil, fmt.Errorf("missing required email configuration (HOST, USERNAME, PASSWORD, FROM_ADDR)")
	}

	return NewEmailService(config, templatesDir)
}

// SendWelcomeEmail sends a welcome email to a new user
func (s *EmailService) SendWelcomeEmail(to, username string) error {
	subject := "Welcome to Mordezzan!"
	templateName := "welcome.html"

	templateData := map[string]interface{}{
		"Username": username,
	}

	return s.SendTemplatedEmail(to, subject, templateName, templateData)
}

// SendTemplatedEmail sends an email using a template
func (s *EmailService) SendTemplatedEmail(to, subject, templateName string, data interface{}) error {
	// Parse template
	var body bytes.Buffer
	if err := s.templates.ExecuteTemplate(&body, templateName, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Prepare email header
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromAddr)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=utf-8"

	// Compose message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	// Setup authentication
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// Setup TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: s.config.InsecureSkipVerify,
		ServerName:         s.config.Host,
	}

	// Connect to server
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to email server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
	}

	// Set the sender and recipient
	if err := client.Mail(s.config.FromAddr); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send the email body
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send data command: %w", err)
	}
	defer wc.Close()

	_, err = wc.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	log.Printf("Email sent to %s with subject: %s", to, subject)
	return nil
}

// SendSimpleEmail sends a simple plain text email
func (s *EmailService) SendSimpleEmail(to, subject, body string) error {
	// Prepare email header
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromAddr)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=utf-8"

	// Compose message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Send the email
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	err := smtp.SendMail(addr, auth, s.config.FromAddr, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Simple email sent to %s with subject: %s", to, subject)
	return nil
}
