package logging

import (
	"io"
	"io/ioutil"
	"os"
)

type Config struct {
	traceHandle   io.Writer
	infoHandle    io.Writer
	warningHandle io.Writer
	errorHandle   io.Writer
}

// DefaultConfig returns default configuration for logger
func DefaultConfig(cfg *Config) *Config {
	if cfg == nil {
		return &Config{
			traceHandle:   ioutil.Discard,
			infoHandle:    os.Stdout,
			warningHandle: os.Stdout,
			errorHandle:   os.Stderr,
		}
	}

	if cfg.traceHandle == nil {
		cfg.traceHandle = ioutil.Discard
	}

	if cfg.infoHandle == nil {
		cfg.infoHandle = os.Stdout
	}

	if cfg.warningHandle == nil {
		cfg.warningHandle = os.Stdout
	}

	if cfg.errorHandle == nil {
		cfg.errorHandle = os.Stderr
	}

	return cfg
}
