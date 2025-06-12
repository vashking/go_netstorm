package attack

import (
	"testing"
	"l7flooder/internal/metrics"
)

func TestStartAttackBasic(t *testing.T) {
	m := &metrics.Metrics{}
	// minimal test: attack on local address, 1 thread, 1 second
	go StartAttack("http://localhost", 1, 1, "", "GET", "", m)
} 