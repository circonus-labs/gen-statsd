// Copyright © 2020 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

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

//Agent is a top level struct for generating and sending StatsD data
type Agent struct {
	id            int
	flushInterval time.Duration
	counterNames  []string
	gaugeNames    []string
	timerNames    []string
	statsdClient  *statsd.Client
}

//CreateAgent creates a new instance of an Agent
func CreateAgent(id int, flush time.Duration, addr, prefix, tags, tagFormat string) (*Agent, error) {

	//Setup some variables
	var client *statsd.Client
	var tagOption statsd.Option
	var tagFormatOption statsd.Option

	//Create the client
	client, err := statsd.New(
		statsd.Address(statsdHost),
		statsd.Network(network),
		statsd.FlushPeriod(flush),
		statsd.Prefix(prefix),
		statsd.ErrorHandler(func(err error) {
			log.Printf("error sending metrics: %s\n", err)
		}),
	)
	if err != nil {
		log.Printf("error creating statsd client: %s\n", err)
	}

	//Check the tagformat
	if tagFormat != "" {
		if tagFormat == "datadog" {
			tagFormatOption = statsd.TagsFormat(statsd.Datadog)
			client = client.Clone(tagFormatOption)
		}
		if tagFormat == "influx" {
			tagFormatOption = statsd.TagsFormat(statsd.InfluxDB)
			client = client.Clone(tagFormatOption)
		}
		return nil, errors.New("unrecognized tag format string")
	}

	//Check for tags
	if tags != "" {
		var err error
		tagOption, err = parseTags(tags)
		if err != nil {
			return nil, err
		}
		client = client.Clone(tagOption)
	}

	a := &Agent{
		id:            id,
		flushInterval: time.Duration(flush),
		counterNames:  genMetricsNames("counter", id, counters),
		gaugeNames:    genMetricsNames("gauge", id, gauges),
		timerNames:    genMetricsNames("timer", id, timers),
		statsdClient:  client,
	}
	return a, nil
}

//Start starts an agent generating and sending metrics to the desired host
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
			log.Printf("flushed %d counters, %d gauges, %d timers for agent %d\n", len(a.counterNames), len(a.gaugeNames), len(a.timerNames), a.id)
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

//flushOnce is to facilitate controlled testing
func (a *Agent) flushOnce() {
	a.statsdClient.Flush()
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
			if tag == "" {
				return nil, errors.New("incomplete key:value pairs")
			}
			tags = append(tags, tag)
		}
	}
	if len(tags) < 2 || len(tags)%2 != 0 {
		return nil, errors.New("incomplete key:value pairs")
	}
	return statsd.Tags(tags...), nil
}
