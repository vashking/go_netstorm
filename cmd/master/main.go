package main

import (
	"flag"
	"l7flooder/internal/cluster"
)

func main() {
	addr := flag.String("addr", ":8080", "Address to listen on")
	flag.Parse()
	master := cluster.NewMaster()
	master.Serve(*addr)
} 