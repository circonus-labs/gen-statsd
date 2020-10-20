package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// config vars, to be manipulated via command line flags
var (
	statsdHost    string
	prefix        string
	flushInterval time.Duration
	spawnInterval time.Duration
	counters      int
	gauges        int
	timers        int
	agents        int
)

func main() {
	flag.StringVar(&statsdHost, "statsd-host", "localhost:8126", "address of statsD host")
	flag.StringVar(&prefix, "prefix", "go-genstatsd", "prefix for metrics")
	flag.DurationVar(&flushInterval, "flush-interval", 10*time.Second, "how often to flush metrics")
	flag.DurationVar(&spawnInterval, "spawn-interval", 10*time.Second, "how often to gen new agents")
	flag.IntVar(&counters, "counters", 50, "number of counters for each agent to hold")
	flag.IntVar(&gauges, "gauges", 30, "number of gauges for each agent to hold")
	flag.IntVar(&timers, "timers", 20, "number of timers for each agent to hold")
	flag.IntVar(&agents, "agents", 10, "max number of agents to run concurrently")
	flag.Parse()

	for i := 0; i < agents; i++ {
		agent := CreateAgent(i, flushInterval, statsdHost, prefix)
		go agent.Start()
		log.Printf("agent %d launched\n", i)
		time.Sleep(spawnInterval)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	<-done
}
