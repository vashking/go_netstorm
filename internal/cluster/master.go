package cluster

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type AttackRequest struct {
	Target   string            `json:"target"`
	Threads  int               `json:"threads"`
	Duration int               `json:"duration"`
	Method   string            `json:"method"`
	Params   map[string]string `json:"params"`
}

type WorkerInfo struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

type Master struct {
	Workers map[string]WorkerInfo
	Mu      sync.Mutex
}

func NewMaster() *Master {
	return &Master{Workers: make(map[string]WorkerInfo)}
}

func (m *Master) RegisterWorker(w http.ResponseWriter, r *http.Request) {
	var info WorkerInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	m.Mu.Lock()
	m.Workers[info.ID] = info
	m.Mu.Unlock()
	w.WriteHeader(200)
}

func (m *Master) ListWorkers(w http.ResponseWriter, r *http.Request) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	json.NewEncoder(w).Encode(m.Workers)
}

func (m *Master) StartAttack(w http.ResponseWriter, r *http.Request) {
	var req AttackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	// gRPC send attack request to workers todo
	w.WriteHeader(202)
}

func (m *Master) Serve(addr string) {
	http.HandleFunc("/register", m.RegisterWorker)
	http.HandleFunc("/workers", m.ListWorkers)
	http.HandleFunc("/attack", m.StartAttack)
	log.Printf("[master] Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
} 