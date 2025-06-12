package main

import (
	"flag"
	"log"
	"l7flooder/internal/cluster"
)

func main() {
	id := flag.String("id", "worker1", "Worker ID")
	master := flag.String("master", "http://localhost:8080", "Master address")
	addr := flag.String("addr", ":9001", "Worker listen address")
	flag.Parse()

	worker := cluster.NewWorker(*id, *master, *addr)
	if err := worker.Register(); err != nil {
		log.Fatalf("Failed to register worker: %v", err)
	}
	worker.ListenAndServe()
} 