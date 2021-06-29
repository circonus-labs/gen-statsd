// Copyright © 2020 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	statsd "gopkg.in/alexcesaro/statsd.v2"
)

func TestParseTags(t *testing.T) {

	// Test setup
	goodTags := "key1:value1,key2:value2,key3:value3"
	badTags1 := "key1:value1,key2:value2,key3"
	badTags2 := "key1:value1,key2:value2,key3:"

	// Run the tests and verify output
	_, err := parseTags(goodTags)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	_, err = parseTags(badTags1)
	if err == nil {
		t.Error("Expected an error but one wasn't returned")
	}

	_, err = parseTags(badTags2)
	if err == nil {
		t.Error("Expected an error but one wasn't returned")
	}
}

func TestGenMetricsNames(t *testing.T) {

	// Test Setup
	id := 0
	metricType := "counter"
	n := 3

	// Run tests and verify output
	metricNames := genMetricsNames(metricType, id, n)
	if len(metricNames) != 3 {
		t.Errorf("Expected length of slice to be %d, got %d", n, len(metricNames))
	}

	for i, name := range metricNames {
		if name != "agent"+strconv.Itoa(id)+"-"+metricType+strconv.Itoa(i) {
			t.Errorf("Expected: %s, got %s", fmt.Sprintf("agent%d-%s%d", id, metricType, i), name)
		}
	}
}

func TestCreateAgent(t *testing.T) {

	srv, err := net.ListenPacket("udp", ":8125")
	if err != nil {
		t.Fatalf("creating server %s", err)
	}
	defer srv.Close()
	go func() {
		buf := make([]byte, 1024)
		_, _, err = srv.ReadFrom(buf)
		if err != nil {
			return
		}
	}()

	// Test Setup

	id := 0
	num := 1
	flush := time.Second * 0
	timerMin := 0
	timerMax := 10
	addr := ":8125"
	prefix := "test"
	goodTags := "key1:value1,key2:value2"
	badTags := "key1:value1,key2:"
	tagFormat1 := "datadog"
	tagFormat2 := "influx"
	badTagFormat := "not a tag format"
	client, err := statsd.New(
		statsd.Address(addr),
		statsd.Network("udp"),
		statsd.FlushPeriod(flush),
		statsd.Prefix(prefix),
		statsd.ErrorHandler(func(err error) {
			log.Printf("error sending metrics to target %s: %s\n", addr, err)
		}),
	)
	if err != nil {
		t.Fatalf("error creating client for target %s: %s", addr, err)
	}

	clients := []*statsd.Client{client}

	// Run tests and verify output
	_, err = CreateAgent(id, num, num, num, timerMax, timerMin, flush, clients, goodTags, tagFormat1, true)
	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}
	_, err = CreateAgent(id, num, num, num, timerMax, timerMin, flush, clients, badTags, tagFormat1, true)
	if err == nil {
		t.Errorf("expected error, got: %s", err)
	}
	_, err = CreateAgent(id, num, num, num, timerMax, timerMin, flush, clients, goodTags, tagFormat2, true)
	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}
	_, err = CreateAgent(id, num, num, num, timerMax, timerMin, flush, clients, goodTags, badTagFormat, true)
	if err == nil {
		t.Errorf("expected error, got: %s", err)
	}
}

func TestNewAgentController(t *testing.T) {
	ac := NewAgentController()
	if ac == nil {
		t.Error("expected agent controller got nil")
	}
}

// func TestAgentControllerStart(t *testing.T) {
// 	ac := NewAgentController()
// 	conf := config{
// 		statsdHost:    ":8125",
// 		network:       "udp",
// 		flushInterval: time.Second * 1,
// 		counters:      1,
// 		gauges:        0,
// 		timers:        0,
// 		agents:        1,
// 		spawnDrift:    10,
// 	}
// 	buf := make([]byte, 1024)
// 	go udpListen(buf)
// 	ac.Start(conf)
// 	time.Sleep(time.Second * 2)
// 	ac.ctx.Done()

// 	fmt.Print(string(buf))
// }

// func udpListen(buf []byte) (int, error) {
// 	srv, err := net.ListenPacket("udp", ":8125")
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer srv.Close()
// 	n, _, err := srv.ReadFrom(buf)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return n, nil
// }
