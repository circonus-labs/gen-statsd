# go-gen-statsd

go-gen-statsd is an experimental tool for load testing with StatsD metrics.

## Usage

```
  Usage of go-gen-statsd:
  -agents int
        max number of agents to run concurrently (default 10)
  -counters int
        number of counters for each agent to hold (default 50)
  -flush-interval duration
        how often to flush metrics (default 10s)
  -gauges int
        number of gauges for each agent to hold (default 30)
  -prefix string
        prefix for metrics (default "go-genstatsd")
  -protocol string
        network protocol to use, tcp or udp (default "udp")
  -spawn-drift int
        spread new agent generation by 0-n seconds (default 10)
  -statsd-host string
        address of statsD host (default "localhost:8125")
  -tag-format string
        format of the tags to send. accepted values "datadog" or "influx" (default "datadog")
  -tags string
        list of K:V comma separated tags. Example: key1:tag1,key2:tag2
  -timers int
        number of timers for each agent to hold (default 20)
```

## Installation

### From Source

1. Clone this repo
1. Navigate to the repo in your filesystem
1. Run `go install`

## Contributing

[Contributing Guide](./CONTRIBUTING.md)