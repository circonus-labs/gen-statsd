package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type Agent struct {
	id            int
	flushInterval time.Duration
	counterNames  []string
	gaugeNames    []string
	timerNames    []string
	statsdClient  *statsd.Client
}

func CreateAgent(id int, flush time.Duration, addr, prefix, tags, tagFormat string) *Agent {
	var tagOption statsd.Option
	var tagFormatOption statsd.Option

	if tagFormat == "datadog" {
		tagFormatOption = statsd.TagsFormat(statsd.Datadog)
	}

	if tagFormat == "influx" {
		tagFormatOption = statsd.TagsFormat(statsd.InfluxDB)
	}

	if tags != "" {
		var err error
		tagOption, err = parseTags(tags)
		if err != nil {
			log.Println(err)
			return nil
		}
	}

	client, err := statsd.New(
		statsd.Address(statsdHost),
		statsd.Network(network),
		statsd.FlushPeriod(flush),
		tagOption,
		tagFormatOption,
		statsd.Prefix(prefix),
		statsd.ErrorHandler(func(err error) {
			log.Printf("error sending metrics: %s\n", err)
		}),
	)
	if err != nil {
		log.Printf("error creating statsd client: %s\n", err)
	}

	a := &Agent{
		id:            id,
		flushInterval: time.Duration(flush),
		counterNames:  genMetricsNames("counter", id, counters),
		gaugeNames:    genMetricsNames("gauge", id, gauges),
		timerNames:    genMetricsNames("timer", id, timers),
		statsdClient:  client,
	}
	return a
}

func (a *Agent) Start(ctx context.Context) {
	ticker := time.NewTicker(a.flushInterval)
	for {
		select {
		case <-ticker.C:
			var wg sync.WaitGroup
			wg.Add(3)
			go func() {
				a.genCounters(ctx)
				wg.Done()
			}()
			go func() {
				a.genGauges(ctx)
				wg.Done()
			}()
			go func() {
				a.genTimers(ctx)
				wg.Done()
			}()
			wg.Wait()
			log.Printf("metrics for agent %d created\n", a.id)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}

}

func (a *Agent) done(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func (a *Agent) genCounters(ctx context.Context) {
	for _, name := range a.counterNames {
		a.statsdClient.Count(name, rand.Intn(10))
		if a.done(ctx) {
			break
		}
	}
}

func (a *Agent) genGauges(ctx context.Context) {
	for _, name := range a.gaugeNames {
		a.statsdClient.Gauge(name, rand.Intn(500))
		if a.done(ctx) {
			break
		}
	}
}

func (a *Agent) genTimers(ctx context.Context) {
	for _, name := range a.timerNames {
		a.statsdClient.Timing(name, rand.Intn(1000))
		if a.done(ctx) {
			break
		}
	}
}

func genMetricsNames(metricType string, id, n int) []string {
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf("agent%d-%s%d", id, metricType, i)
	}
	return names
}

func parseTags(t string) (statsd.Option, error) {
	kvp := strings.Split(t, ",")
	var tags []string
	for _, pairs := range kvp {
		kv := strings.Split(pairs, ":")
		for _, tag := range kv {
			tags = append(tags, tag)
		}
	}
	if len(tags)%2 != 0 {
		return nil, errors.New("incomplete key:value pairs")
	}
	return statsd.Tags(tags...), nil
}
