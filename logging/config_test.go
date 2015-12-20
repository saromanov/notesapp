package logging

import (
  "testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig(nil)
	if cfg.traceHandle == nil {
		t.Errorf("Must be non-empty")
	}

	if cfg.infoHandle == nil {
		t.Errorf("Must be non-empty")
	}

	if cfg.warningHandle == nil {
		t.Errorf("Must be non-empty")
	}

	if cfg.errorHandle == nil {
		t.Errorf("Must be non-empty")
	}
}