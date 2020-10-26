package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/sys/unix"
)

// config vars, to be manipulated via command line flags
var (
	statsdHost    string
	prefix        string
	network       string
	tags          string
	tagFormat     string
	flushInterval time.Duration
	spawnDrift    int
	counters      int
	gauges        int
	timers        int
	agents        int
)

func main() {
	defaultPrefix, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&statsdHost, "statsd-host", "localhost:8125", "address of statsD host")
	flag.StringVar(&prefix, "prefix", filepath.Base(defaultPrefix), "prefix for metrics")
	flag.DurationVar(&flushInterval, "flush-interval", 10*time.Second, "how often to flush metrics")
	flag.IntVar(&spawnDrift, "spawn-drift", 10, "spread new agent generation by 0-n seconds")
	flag.StringVar(&network, "protocol", "udp", "network protocol to use, tcp or udp")
	flag.StringVar(&tagFormat, "tag-format", "datadog", "format of the tags to send. accepted values \"datadog\" or \"influx\"")
	flag.StringVar(&tags, "tags", "", "list of K:V comma separated tags. Example: key1:tag1,key2:tag2")
	flag.IntVar(&counters, "counters", 50, "number of counters for each agent to hold")
	flag.IntVar(&gauges, "gauges", 30, "number of gauges for each agent to hold")
	flag.IntVar(&timers, "timers", 20, "number of timers for each agent to hold")
	flag.IntVar(&agents, "agents", 10, "max number of agents to run concurrently")
	flag.Parse()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, unix.SIGTERM, unix.SIGHUP, unix.SIGPIPE, unix.SIGINFO)
	go func() {
		sig := <-sigs
		log.Printf("received %s, exiting\n", sig.String())
		signal.Stop(sigs)
		signal.Reset() // so a second ctrl-c will force immediate stop (if user in hurry and a go routine is in a timer)
		cancel()
	}()

	for i := 0; i < agents; i++ {
		wg.Add(1)
		go func(id int) {
			agent, err := CreateAgent(id, flushInterval, statsdHost, prefix, tags, tagFormat)
			if err != nil {
				log.Printf("error instantiating agent%d: %s", id, err)
				os.Exit(1)
			}
			log.Printf("launching agent %d\n", id)
			agent.Start(ctx)
			wg.Done()
		}(i)
		time.Sleep(time.Duration(rand.Intn(spawnDrift)) * time.Second)
		if done(ctx) {
			break
		}
	}

	wg.Wait()
}

func done(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
