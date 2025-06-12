package proxy

import (
	"bufio"
	"os"
	"sync"
)

type ProxyManager struct {
	proxies []string
	mu      sync.Mutex
	index   int
}

func NewProxyManager(path string) (*ProxyManager, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			proxies = append(proxies, line)
		}
	}
	return &ProxyManager{proxies: proxies}, nil
}

func (pm *ProxyManager) Next() string {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if len(pm.proxies) == 0 {
		return ""
	}
	proxy := pm.proxies[pm.index]
	pm.index = (pm.index + 1) % len(pm.proxies)
	return proxy
}
