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
	"os"
	"strings"
	"sync"
	"time"

	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// AgentController is the main controller of the agents
type AgentController struct { //nolint:govet
	sig  chan os.Signal
	wg   sync.WaitGroup
	ctx  context.Context
	cncl context.CancelFunc
}

// Agent is a struct for generating and sending StatsD data
type Agent struct {
	counterNames  []string
	gaugeNames    []string
	timerNames    []string
	statsdClients []*statsd.Client
	id            int
	flushInterval time.Duration
	timerSamples  int
	valueMax      int
	valueMin      int
	quiet         bool
}

// NewAgentController creates a new agent pool
func NewAgentController() *AgentController {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithCancel(context.Background())
	return &AgentController{
		sig:  make(chan os.Signal, 1),
		wg:   sync.WaitGroup{},
		ctx:  ctx,
		cncl: cancel,
	}
}

// Start kicks off the main process of sending statsD metrics from agents
func (ac *AgentController) Start(c config) {

	if c.runTime > time.Duration(0) {
		go func() {
			time.Sleep(c.runTime)
			log.Printf("run time (%s) reached, canceling", c.runTime.String())
			ac.cncl()
		}()
	}

	ac.signalNotifySetup()
	go ac.handleSignals()

	targets := strings.Split(c.statsdHosts, ",")
	statsdClients := make([]*statsd.Client, 0)
	for _, t := range targets {
		t := t
		proto := "udp"
		address := ":8125"
		spec := strings.Split(t, ":")
		switch len(spec) {
		case 3:
			address = fmt.Sprintf("%s:%s", spec[0], spec[1])
			proto = spec[2]
		case 2:
			address = fmt.Sprintf("%s:%s", spec[0], spec[1])
		case 1:
			address = spec[0]
			if isSocket(address) {
				proto = "unixgram"
			}
		default:
			log.Printf("invalid target spec (%s)", t)
			continue
		}
		client, err := statsd.New(
			statsd.Address(address),
			statsd.Network(proto),
			statsd.FlushPeriod(c.flushInterval),
			statsd.Prefix(c.prefix),
			statsd.ErrorHandler(func(err error) {
				log.Printf("error sending metrics to target %s: %s\n", t, err)
			}),
		)
		if err != nil {
			log.Printf("error creating client for target %s: %s", t, err)
			continue
		}
		if c.sampleRate > 0 {
			client = client.Clone(statsd.SampleRate(float32(c.sampleRate)))
		}
		statsdClients = append(statsdClients, client)
	}

	if len(statsdClients) == 0 {
		log.Fatal("no targets defined")
	}

	for i := 0; i < c.agents; i++ {
		ac.wg.Add(1)
		go func(id int) {
			agent, err := CreateAgent(id, c.counters, c.gauges, c.timers, c.valueMax, c.valueMin, c.flushInterval, statsdClients, c.tags, c.tagFormat, c.quiet)
			if err != nil {
				log.Printf("error instantiating agent%d: %s\n", id, err)
				ac.ctx.Done()
				ac.wg.Done()
				return
			}
			agent.timerSamples = c.timerSamples
			if !c.quiet {
				log.Printf("launching agent %d\n", id)
			}
			agent.Start(ac.ctx)
			ac.wg.Done()
		}(i)
		if done(ac.ctx) {
			break
		}
		time.Sleep(time.Duration(rand.Intn(c.spawnDrift)) * time.Second) //nolint:gosec
	}

	ac.wg.Wait()
}

// CreateAgent creates a new instance of an Agent
func CreateAgent(id, counters, gauges, timers, max, min int, flush time.Duration, targets []*statsd.Client, tags, tagFormat string, quiet bool) (*Agent, error) {

	// Check the tagformat
	if tagFormat != "" {
		for idx, c := range targets {
			client, err := parseTagFormat(c, tagFormat)
			if err != nil {
				return nil, err
			}
			targets[idx] = client
		}
	}

	// Check for tags
	if tags != "" {
		tagOption, err := parseTags(tags)
		if err != nil {
			return nil, err
		}
		for idx, c := range targets {
			targets[idx] = c.Clone(tagOption)
		}
	}

	a := &Agent{
		id:            id,
		flushInterval: time.Duration(flush),
		counterNames:  genMetricsNames("counter", id, counters),
		gaugeNames:    genMetricsNames("gauge", id, gauges),
		timerNames:    genMetricsNames("timer", id, timers),
		valueMax:      max,
		valueMin:      min,
		statsdClients: targets,
		quiet:         quiet,
	}
	return a, nil
}

// Start starts an agent generating and sending metrics to the desired host
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
			if !a.quiet {
				log.Printf("flushed %d counters, %d gauges, %d timers(*%d samples) for agent %d\n", len(a.counterNames), len(a.gaugeNames), len(a.timerNames), a.timerSamples, a.id)
			}
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

func (a *Agent) getSample() int {
	if a.valueMax == a.valueMin {
		return a.valueMax
	}
	return rand.Intn(a.valueMax-a.valueMin+1) + a.valueMin //nolint:gosec
}

func (a *Agent) genCounters(ctx context.Context) {
	for _, name := range a.counterNames {
		val := a.getSample()
		for _, c := range a.statsdClients {
			c.Count(name, val)
			if a.done(ctx) {
				break
			}
		}
		if a.done(ctx) {
			break
		}
	}
}

func (a *Agent) genGauges(ctx context.Context) {
	for _, name := range a.gaugeNames {
		val := a.getSample()
		for _, c := range a.statsdClients {
			c.Gauge(name, val)
			if a.done(ctx) {
				break
			}
		}
		if a.done(ctx) {
			break
		}
	}
}

func (a *Agent) genTimers(ctx context.Context) {
	for _, name := range a.timerNames {
		for i := 0; i < a.timerSamples; i++ {
			val := a.getSample()
			for _, c := range a.statsdClients {
				c.Timing(name, val)
				if a.done(ctx) {
					break
				}
			}
			if a.done(ctx) {
				break
			}
		}
		if a.done(ctx) {
			break
		}
	}
}

// flushOnce is to facilitate controlled testing
func (a *Agent) flushOnce() { //nolint:go-lint,unused
	for _, c := range a.statsdClients {
		c.Flush()
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

func parseTagFormat(client *statsd.Client, tf string) (*statsd.Client, error) {
	if tf == "datadog" {
		return client.Clone(statsd.TagsFormat(statsd.Datadog)), nil
	}
	if tf == "influx" {
		return client.Clone(statsd.TagsFormat(statsd.InfluxDB)), nil
	}
	return nil, errors.New("unrecognized tag format")
}

func done(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func isSocket(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSocket != 0
}
