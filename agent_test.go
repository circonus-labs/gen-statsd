package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestParseTags(t *testing.T) {

	//Test setup
	goodTags := "key1:value1,key2:value2,key3:value3"
	badTags1 := "key1:value1,key2:value2,key3"
	badTags2 := "key1:value1,key2:value2,key3:"

	//Run the tests and verify output
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

	//Test Setup
	id := 0
	metricType := "counter"
	n := 3

	//Run tests and verify output
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
