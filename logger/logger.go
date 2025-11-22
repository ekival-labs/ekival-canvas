package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ANSI color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// isTerminal checks if the output is a terminal (for color support)
func isTerminal(w io.Writer) bool {
	// Check if we are in a test environment
	isTest := isTestEnvironment() // Helper function to determine if we are in a test
	if isTest {
		return true // Always assume terminal for tests to force color output
	}
	if f, ok := w.(*os.File); ok {
		// On Windows, check if it's a console
		if runtime.GOOS == "windows" {
			// Try to enable ANSI color support on Windows 10+
			// This is a best-effort approach
			return true // Assume terminal for tests
		}
		// For Unix-like systems, check if it's a terminal
		stat, err := f.Stat()
		if err != nil {
			return false
		}
		return (stat.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// isTestEnvironment checks if the current execution is within a Go test.
func isTestEnvironment() bool {
	argsStr := strings.Join(os.Args, " ")
	exeName := os.Args[0]

	if strings.HasSuffix(exeName, ".test") ||
		strings.HasSuffix(exeName, "_test.exe") ||
		strings.Contains(exeName, "_test") ||
		strings.Contains(argsStr, "-test.") ||
		strings.Contains(argsStr, "go test") ||
		strings.Contains(argsStr, "-test.v") ||
		strings.Contains(argsStr, "-test.run") ||
		strings.Contains(argsStr, "-run") ||
		strings.Contains(argsStr, "-v") {
		return true
	}

	if tmpdir := os.Getenv("TMPDIR"); tmpdir != "" && strings.Contains(tmpdir, "go-build") {
		return true
	}
	if tmp := os.Getenv("TMP"); tmp != "" && strings.Contains(tmp, "go-build") {
		return true
	}

	if os.Getenv("ZEROLOG_CONSOLE") == "true" {
		return true
	}

	if os.Getenv("GO_TEST") != "" {
		return true
	}
	return false
}

// enableWindowsColors attempts to enable ANSI color support on Windows
func enableWindowsColors() {
	if runtime.GOOS == "windows" {
		// Try to enable virtual terminal processing for Windows 10+
		// This is done via kernel32.dll, but we'll rely on the terminal
		// being configured correctly. Most modern terminals support ANSI.
		// For PowerShell and modern cmd.exe, colors should work.
	}
}

// init automatically initializes the global zerolog logger with colorful console output
// This ensures that even in test environments, logs are displayed in a readable format
func init() {
	// Check if we're in a test environment using the helper function
	isTest := isTestEnvironment()

	// Always initialize with console format if in test mode or if explicitly requested
	shouldInit := isTest || os.Getenv("ZEROLOG_CONSOLE") == "true"

	if shouldInit {
		// Enable Windows colors if on Windows
		enableWindowsColors()

		// Use colorable output for Windows compatibility
		var output io.Writer = os.Stdout
		if runtime.GOOS == "windows" {
			output = colorable.NewColorableStdout()
		}

		// Create a colorful console writer for stdout
		// Force colors on for tests, even if not detected as TTY
		consoleWriter := zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
			NoColor:    false, // Force colors on for tests
			FormatLevel: func(i interface{}) string {
				if ll, ok := i.(string); ok {
					switch ll {
					case "trace":
						return fmt.Sprintf("%s%s%s", colorGray, ll, colorReset)
					case "debug":
						return fmt.Sprintf("%s%s%s", colorCyan, ll, colorReset)
					case "info":
						return fmt.Sprintf("%s%s%s%s", colorBold, colorGreen, ll, colorReset)
					case "warn":
						return fmt.Sprintf("%s%s%s%s", colorBold, colorYellow, ll, colorReset)
					case "error":
						return fmt.Sprintf("%s%s%s%s", colorBold, colorRed, ll, colorReset)
					case "fatal":
						return fmt.Sprintf("%s%s%s%s%s", colorBold, colorRed, colorBold, ll, colorReset)
					case "panic":
						return fmt.Sprintf("%s%s%s%s%s", colorBold, colorRed, colorBold, ll, colorReset)
					default:
						return fmt.Sprintf("%s%s%s", colorWhite, ll, colorReset)
					}
				}
				return ""
			},
			FormatFieldName: func(i interface{}) string {
				return fmt.Sprintf("%s%s%s", colorCyan, i, colorReset)
			},
			FormatFieldValue: func(i interface{}) string {
				return fmt.Sprintf("%s%v%s", colorYellow, i, colorReset)
			},
			FormatCaller: func(i interface{}) string {
				return fmt.Sprintf("%s%s%s", colorGray, i, colorReset)
			},
			FormatTimestamp: func(i interface{}) string {
				return fmt.Sprintf("%s%s%s", colorBlue, i, colorReset)
			},
		}

		// Set the global logger to use console format with colors
		zerolog.TimeFieldFormat = time.RFC3339
		// Set log level to Debug for tests to see more information
		logLevel := zerolog.DebugLevel
		log.Logger = zerolog.New(consoleWriter).Level(logLevel).With().Timestamp().Caller().Logger()
	}
}

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

	// Use colorable output for Windows compatibility if writing to stdout
	var output io.Writer = consoleOutputWriter
	if runtime.GOOS == "windows" {
		if consoleOutputWriter == os.Stdout {
			output = colorable.NewColorableStdout()
		} else if consoleOutputWriter == os.Stderr {
			output = colorable.NewColorableStderr()
		}
	}

	// Console writer for stdout or custom writer with colorful output
	consoleWriter := zerolog.ConsoleWriter{
		Out:        output,
		TimeFormat: time.RFC3339,
		NoColor:    false,
		FormatLevel: func(i interface{}) string {
			if ll, ok := i.(string); ok {
				switch ll {
				case "trace":
					return fmt.Sprintf("%s%s%s", colorGray, ll, colorReset)
				case "debug":
					return fmt.Sprintf("%s%s%s", colorCyan, ll, colorReset)
				case "info":
					return fmt.Sprintf("%s%s%s%s", colorBold, colorGreen, ll, colorReset)
				case "warn":
					return fmt.Sprintf("%s%s%s%s", colorBold, colorYellow, ll, colorReset)
				case "error":
					return fmt.Sprintf("%s%s%s%s", colorBold, colorRed, ll, colorReset)
				case "fatal":
					return fmt.Sprintf("%s%s%s%s%s", colorBold, colorRed, colorBold, ll, colorReset)
				case "panic":
					return fmt.Sprintf("%s%s%s%s%s", colorBold, colorRed, colorBold, ll, colorReset)
				default:
					return fmt.Sprintf("%s%s%s", colorWhite, ll, colorReset)
				}
			}
			return ""
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s%s%s", colorCyan, i, colorReset)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s%v%s", colorYellow, i, colorReset)
		},
		FormatCaller: func(i interface{}) string {
			return fmt.Sprintf("%s%s%s", colorGray, i, colorReset)
		},
		FormatTimestamp: func(i interface{}) string {
			return fmt.Sprintf("%s%s%s", colorBlue, i, colorReset)
		},
	}
	writers = append(writers, CustomLevelWriter{Writer: consoleWriter, Level: consoleLogLevel})
	log.Debug().Str("console_log_level", consoleLogLevel.String()).Msg("Console writer added with specific level.")

	// File writer with rotation
	if logFilePath != "" {
		log.Debug().Str("log_file_path", logFilePath).Str("file_log_level", fileLogLevel.String()).Msg("Configuring file rotation for logs.")
		logWriter, err := rotatelogs.New(
			logFilePath+".%Y%m%d",
			rotatelogs.WithLinkName(logFilePath),
			rotatelogs.WithMaxAge(30*24*time.Hour),    // Keep logs for 30 days
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
