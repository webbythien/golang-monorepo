package l

import (
	"os"
	"testing"
)

func TestL(t *testing.T) {
	os.Setenv("LOG_LEVEL", "WARN")
	ll := New()
	ll.Info("test info")
	ll.Debug("test debug")
	ll.Error("test error")
	ll.Warn("test warn")
}

func TestL2(t *testing.T) {
	os.Setenv("LOG_ENCODER", "loki")
	os.Setenv("LOG_LEVEL", "DEBUG")
	ll := New()
	ll.Info("test info", String("test", "test"))
	ll.Debug("test debug")
	ll.Error("test error")
	ll.Warn("test warn")
}
