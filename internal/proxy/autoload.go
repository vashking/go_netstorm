package proxy

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type AutoProxyManager struct {
	proxies []string
	mu      sync.Mutex
	index   int
}

func NewAutoProxyManager(source string) (*AutoProxyManager, error) {
	var proxies []string
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		resp, err := http.Get(source)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		proxies = parseProxies(resp.Body)
	} else {
		file, err := os.Open(source)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		proxies = parseProxies(file)
	}
	return &AutoProxyManager{proxies: proxies}, nil
}

func (apm *AutoProxyManager) Next() string {
	apm.mu.Lock()
	defer apm.mu.Unlock()
	if len(apm.proxies) == 0 {
		return ""
	}
	proxy := apm.proxies[apm.index]
	apm.index = (apm.index + 1) % len(apm.proxies)
	return proxy
}

func (apm *AutoProxyManager) Reload(source string) error {
	var proxies []string
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		resp, err := http.Get(source)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		proxies = parseProxies(resp.Body)
	} else {
		file, err := os.Open(source)
		if err != nil {
			return err
		}
		defer file.Close()
		proxies = parseProxies(file)
	}
	apm.mu.Lock()
	apm.proxies = proxies
	apm.index = 0
	apm.mu.Unlock()
	return nil
}

func parseProxies(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	var proxies []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			proxies = append(proxies, line)
		}
	}
	return proxies
} 