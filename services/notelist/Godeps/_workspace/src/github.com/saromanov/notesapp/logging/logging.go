package logging

import (
	"log"
)

// Logger provides
type Logger struct {
	trace   *log.Logger
	info    *log.Logger
	warning *log.Logger
	errorm  *log.Logger
}

// NewLOgger provides configuration for logger object
func NewLogger(config *Config) *Logger {
	cfg := DefaultConfig(config)
	logger := new(Logger)
	logger.trace = log.New(cfg.traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime)

	logger.info = log.New(cfg.infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime)

	logger.warning = log.New(cfg.warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime)

	logger.errorm = log.New(cfg.errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime)

	return logger
}

func (logger *Logger) Info(msg string) {
	logger.info.Println(msg)
}

func (logger *Logger) Trace(msg string) {
	logger.trace.Println(msg)
}

func (logger *Logger) Warning(msg string) {
	logger.warning.Println(msg)
}

func (logger *Logger) Error(msg string) {
	logger.errorm.Println(msg)
}
