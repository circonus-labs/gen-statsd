package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type Agent struct {
	ID            int
	FlushInterval time.Duration
	CounterNames  []string
	GaugeNames    []string
	TimerNames    []string
	StatsdClient  *statsd.Client
}

func (a *Agent) Start() {
	for {
		go genCounters(a.CounterNames, a.StatsdClient)
		go genGauges(a.GaugeNames, a.StatsdClient)
		go genTimers(a.TimerNames, a.StatsdClient)
		log.Printf("metrics for agent %d created\n", a.ID)
		time.Sleep(a.FlushInterval)
	}
}

func CreateAgent(id int, flush time.Duration, addr, prefix string) *Agent {
	client, err := statsd.New(
		statsd.Address(statsdHost),
		statsd.FlushPeriod(flush),
		statsd.Prefix(prefix),
		statsd.ErrorHandler(func(err error) {
			log.Printf("error sending metrics: %d\n", err)
		}),
	)
	if err != nil {
		log.Printf("error creating statsd client: %s\n", err)
	}

	a := &Agent{
		ID:            id,
		FlushInterval: time.Duration(flush),
		CounterNames:  genMetricsNames("counter", id, counters),
		GaugeNames:    genMetricsNames("gauge", id, gauges),
		TimerNames:    genMetricsNames("timer", id, timers),
		StatsdClient:  client,
	}
	return a
}

func genMetricsNames(metricType string, id, n int) []string {
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf("%s-agent%d-%s%d", prefix, id, metricType, i)
	}
	return names
}

func genCounters(names []string, client *statsd.Client) {
	for _, name := range names {
		client.Count(name, rand.Intn(10))
	}
}

func genGauges(names []string, client *statsd.Client) {
	for _, name := range names {
		client.Gauge(name, rand.Intn(500))
	}
}

func genTimers(names []string, client *statsd.Client) {
	for _, name := range names {
		client.Timing(name, rand.Intn(1000))
	}
}
