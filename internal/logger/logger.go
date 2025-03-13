package logger

import (
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.SugaredLogger
)

// Config holds logger configuration
type Config struct {
	LogLevel         LogLevel
	IncludeTimestamp bool
	IncludeFileLine  bool
	Development      bool
	Output           interface{} // Can be os.Stdout, os.Stderr, io.Writer, etc.
}

// LogLevel defines logging severity
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

// Init initializes the logger with the given configuration
func Init(config Config) {
	// Convert our custom log level to Zap's log level
	var level zapcore.Level
	switch config.LogLevel {
	case LogLevelDebug:
		level = zapcore.DebugLevel
	case LogLevelInfo:
		level = zapcore.InfoLevel
	case LogLevelWarning:
		level = zapcore.WarnLevel
	case LogLevelError:
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Configure zap logger
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      config.Development,
		Encoding:         "console", // could also be "json" for production environments
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// If a custom output is specified, we'll handle it after building the logger

	// Customize encoder config
	zapConfig.EncoderConfig.TimeKey = "time"
	zapConfig.EncoderConfig.LevelKey = "level"
	zapConfig.EncoderConfig.NameKey = "logger"
	zapConfig.EncoderConfig.CallerKey = "caller"
	zapConfig.EncoderConfig.MessageKey = "msg"
	zapConfig.EncoderConfig.StacktraceKey = "stacktrace"

	// Time format
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Configure caller display (file:line)
	if config.IncludeFileLine {
		zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	} else {
		// Don't include caller info
		zapConfig.EncoderConfig.CallerKey = ""
	}

	// Create the logger
	logger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	// If a custom output was provided (used mainly for testing)
	if config.Output != nil {
		// We need to create a custom core
		if writer, ok := config.Output.(io.Writer); ok {
			// Create a custom zapcore with the provided writer
			encoderConfig := zapConfig.EncoderConfig
			encoder := zapcore.NewConsoleEncoder(encoderConfig)
			writeSyncer := zapcore.AddSync(writer)
			core := zapcore.NewCore(encoder, writeSyncer, zap.NewAtomicLevelAt(level))

			// Create a new logger with the custom core
			logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		}
	}

	log = logger.Sugar()
}

// Debug logs a debug-level message
func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info logs an info-level message
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warning logs a warning-level message
func Warning(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error logs an error-level message
func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// ErrorWithStack logs an error with a stack trace
func ErrorWithStack(err error) {
	log.With("error", err).Error("Error with stack trace")
}

// Fatal logs a fatal-level message and then exits
func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// LogRequest logs HTTP request information
func LogRequest(r *http.Request, duration time.Duration) {
	log.Infow("HTTP Request",
		"method", r.Method,
		"uri", r.RequestURI,
		"remoteAddr", r.RemoteAddr,
		"duration", duration,
	)
}

// With returns a logger with the given key-value pairs
func With(args ...interface{}) *zap.SugaredLogger {
	return log.With(args...)
}

// GetZapLogger returns the underlying zap logger for advanced usage
func GetZapLogger() *zap.Logger {
	return log.Desugar()
}
