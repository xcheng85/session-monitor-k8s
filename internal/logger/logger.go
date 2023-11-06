package logger

// level-based logging
// Using the zap config struct to create a logger
type Level string

// inversion of dependencies principle,
// caller satisfy the struct
type LogConfig struct {
	Environment string
	LogLevel    Level
}

const (
	FATAL Level = "FATAL"
	DEBUG Level = "DEBUG"
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
	PANIC Level = "PANIC"
)
