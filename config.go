// Copyright © 2020 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/namsral/flag"
)

type config struct {
	statsdHosts   string
	prefix        string
	tags          string
	tagFormat     string
	flushInterval time.Duration
	spawnDrift    int
	counters      int
	gauges        int
	timers        int
	timerValueMax int
	timerValueMin int
	tsamples      int
	agents        int
	version       bool
}

func genConfig() config {
	c := config{}

	defaultPrefix, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	flag.StringVar(&c.statsdHosts, "statsd-host", "localhost:8125:udp", "comma separated list of ip:port:proto for statsD host(s)")
	flag.StringVar(&c.prefix, "prefix", filepath.Base(defaultPrefix), "prefix for metrics")
	flag.DurationVar(&c.flushInterval, "flush-interval", 10*time.Second, "how often to flush metrics")
	flag.IntVar(&c.spawnDrift, "spawn-drift", 10, "spread new agent generation by 0-n seconds")
	flag.StringVar(&c.tagFormat, "tag-format", "", "format of the tags to send. accepted values \"datadog\" or \"influx\"")
	flag.StringVar(&c.tags, "tags", "", "list of K:V comma separated tags. Example: key1:tag1,key2:tag2")
	flag.IntVar(&c.counters, "counters", 50, "number of counters for each agent to hold")
	flag.IntVar(&c.timerValueMax, "timer-value-max", 100, "max timer value")
	flag.IntVar(&c.timerValueMin, "timer-value-min", 0, "min timer value")
	flag.IntVar(&c.gauges, "gauges", 30, "number of gauges for each agent to hold")
	flag.IntVar(&c.timers, "timers", 20, "number of timers for each agent to hold")
	flag.IntVar(&c.tsamples, "timer-samples", 10, "number of timer samples per iteration")
	flag.IntVar(&c.agents, "agents", 10, "max number of agents to run concurrently")
	flag.BoolVar(&c.version, "version", false, "show version information")
	flag.Parse()

	return c
}
