package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// CustomLevelWriter wraps an io.Writer and filters logs based on a minimum level.
type CustomLevelWriter struct {
	Writer io.Writer
	Level  zerolog.Level
}

// WriteLevel implements the zerolog.LevelWriter interface.
func (lw CustomLevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= lw.Level {
		return lw.Writer.Write(p)
	}
	return len(p), nil
}

// Write implements the io.Writer interface.
func (lw CustomLevelWriter) Write(p []byte) (n int, err error) {
	// Default to info level if Write is called directly without WriteLevel
	return lw.WriteLevel(zerolog.InfoLevel, p)
}

// Logger holds the zerolog.Logger instance.
type Logger struct {
	*zerolog.Logger
}

// NewLogger initializes and returns a new Logger instance.
// It configures zerolog for JSON output to stdout and optionally to a file with daily rotation.
// It also sets up human-readable console output for development environments.
func NewLogger(logFilePath string, consoleLogLevel zerolog.Level, fileLogLevel zerolog.Level) *Logger {
	return newLoggerInternal(logFilePath, consoleLogLevel, fileLogLevel, os.Stdout)
}

// NewTestLogger initializes and returns a new Logger instance for testing.
// It allows specifying a custom writer for console output (e.g., a bytes.Buffer).
func NewTestLogger(logFilePath string, consoleLogLevel zerolog.Level, fileLogLevel zerolog.Level, consoleOutputWriter io.Writer) *Logger {
	return newLoggerInternal(logFilePath, consoleLogLevel, fileLogLevel, consoleOutputWriter)
}

// newLoggerInternal is the internal implementation for creating a Logger.
func newLoggerInternal(logFilePath string, consoleLogLevel zerolog.Level, fileLogLevel zerolog.Level, consoleOutputWriter io.Writer) *Logger {
	log.Trace().Str("log_file_path", logFilePath).
		Str("console_log_level", consoleLogLevel.String()).
		Str("file_log_level", fileLogLevel.String()).
		Msg("Starting NewLogger function.")

	var writers []io.Writer

	// Console writer for stdout or custom writer
	consoleWriter := zerolog.ConsoleWriter{Out: consoleOutputWriter, TimeFormat: time.RFC3339}
	writers = append(writers, CustomLevelWriter{Writer: consoleWriter, Level: consoleLogLevel})
	log.Debug().Str("console_log_level", consoleLogLevel.String()).Msg("Console writer added with specific level.")

	// File writer with rotation
	if logFilePath != "" {
		log.Debug().Str("log_file_path", logFilePath).Str("file_log_level", fileLogLevel.String()).Msg("Configuring file rotation for logs.")
		logWriter, err := rotatelogs.New(
			logFilePath+".%Y%m%d",
			rotatelogs.WithLinkName(logFilePath),
			rotatelogs.WithMaxAge(30*24*time.Hour),     // Keep logs for 30 days
			rotatelogs.WithRotationTime(24*time.Hour), // Rotate daily
		)
		if err != nil {
			log.Fatal().Err(err).Str("log_file_path", logFilePath).Msg("Failed to create rotatelogs writer, application cannot start.")
		}
		writers = append(writers, CustomLevelWriter{Writer: logWriter, Level: fileLogLevel})
		log.Debug().Str("log_file_path", logFilePath).Str("file_log_level", fileLogLevel.String()).Msg("File writer with rotation added.")
	} else {
		log.Debug().Msg("No log file path provided, skipping file logging.")
	}

	multiWriter := zerolog.MultiLevelWriter(writers...)
	log.Debug().Msg("Multi-level writer created for zerolog.")

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return fmt.Sprintf("%+v", err)
	}
	log.Debug().Msg("Global zerolog settings configured (TimeFormat, ErrorStackMarshaler).")

	// Determine the minimum level to ensure all relevant messages are processed by the multi-level writer
	minLevel := consoleLogLevel
	if fileLogLevel < minLevel {
		minLevel = fileLogLevel
	}

	// Integrate AddSource functionality using zerolog.With().Caller()
	zlog := zerolog.New(multiWriter).Level(minLevel).With().Timestamp().Caller().Logger()
	log.Debug().Str("min_effective_log_level", minLevel.String()).Msg("Zerolog instance created with timestamp and caller (source).")

	log.Logger = zlog // Set the global zerolog logger to our configured instance
	log.Info().Msg("Global zerolog logger initialized successfully.")
	log.Trace().Msg("Finished NewLogger function.")

	return &Logger{&zlog}
}

// WithContext returns a new logger with additional context fields.
func (l *Logger) WithContext(fields map[string]interface{}) *Logger {
	log.Trace().Interface("fields", fields).Msg("Starting WithContext method.")
	ctx := l.Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	newLogger := ctx.Logger()
	log.Trace().Msg("Finished WithContext method, new logger with context created.")
	return &Logger{&newLogger}
}

// LogError logs an error with stack trace.
func (l *Logger) LogError(err error, source string, fields map[string]interface{}) {
	log.Trace().Err(err).Str("source", source).Interface("fields", fields).Msg("Starting LogError method.")
	event := l.Error().Stack().Err(errors.WithStack(err)).Str("source", source)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg("Error occurred")
	log.Trace().Msg("Finished LogError method.")
}

// LogPanic logs a panic with stack trace.
func (l *Logger) LogPanic(err interface{}, stack []byte, source string, fields map[string]interface{}) {
	log.Trace().Interface("panic_value", err).Bytes("stack_trace_bytes", stack).Str("source", source).Interface("fields", fields).Msg("Starting LogPanic method.")
	event := l.Error().Stack().
		Interface("panic", err).
		Bytes("stack_trace", stack).
		Str("source", source)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg("Panic occurred")
	log.Trace().Msg("Finished LogPanic method.")
}

// StartSilentLogger initializes and returns a new SilentLogger instance.
// It configures zerolog to discard all log events, effectively making it silent.
func StartSilentLogger() *Logger {
	// Create a zerolog.Logger that discards all events
	nopLogger := zerolog.Nop()
	return &Logger{&nopLogger}
}

// NewLoggerFromZerolog creates a new Logger instance from an existing zerolog.Logger.
func NewLoggerFromZerolog(zl *zerolog.Logger) *Logger {
	return &Logger{zl}
}
