package metrics

import (
	"fmt"
	"sync"
	"time"
)

type Metrics struct {
	mu        sync.Mutex
	Requests  int
	Success   int
	Fail      int
	Latency   time.Duration
}

func (m *Metrics) Add(success bool, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Requests++
	if success {
		m.Success++
	} else {
		m.Fail++
	}
	m.Latency += latency
}

func (m *Metrics) Print() {
	m.mu.Lock()
	defer m.mu.Unlock()
	avgLatency := time.Duration(0)
	if m.Requests > 0 {
		avgLatency = m.Latency / time.Duration(m.Requests)
	}
	fmt.Printf("Requests: %d | Success: %d | Fail: %d | Avg Latency: %v\n", m.Requests, m.Success, m.Fail, avgLatency)
}
