# gen-statsd

[![Go Report Card](https://goreportcard.com/badge/github.com/circonus-labs/gen-statsd)](https://goreportcard.com/report/github.com/circonus-labs/gen-statsd)

gen-statsd is an experimental tool written in Go for load and signal testing with StatsD metrics.

## Configuration

### Command-line Flags

```
Usage of ./gen-statsd:
  -agents=10: max number of agents to run concurrently
  -config="": path to config file
  -counters=50: number of counters for each agent to hold
  -flush-interval=10s: how often to flush metrics
  -gauges=30: number of gauges for each agent to hold
  -prefix="gen-statsd": prefix for metrics
  -run-time=0s: how long to run, 0=forever
  -sample-rate=0: sampling rate
  -spawn-drift=10: spread new agent generation by 0-n seconds
  -statsd-hosts="localhost:8125:udp": comma separated list of ip:port:proto for statsD host(s)
  -tag-format="": format of the tags to send. accepted values "datadog" or "influx"
  -tags="": list of K:V comma separated tags. Example: key1:tag1,key2:tag2
  -timer-samples=10: number of timer samples per iteration
  -timers=20: number of timers for each agent to hold
  -value-max=100: maximum value
  -value-min=0: minimum value
  -version=false: show version information
```

### Environment Variables

  | Variable      |                       Description                               |
  |:-------------:|:----------------------------------------------------------------|
  |AGENTS         |max number of agents to run concurrently (default 10)            |
  |CONFIG         |path to config file                                              |
  |COUNTERS       |number of counters for each agent to hold (default 50)           |
  |FLUSH_INTERVAL |how often to flush metrics (default 10s)                         |
  |GAUGES         |number of gauges for each agent to hold (default 30)             |
  |PREFIX         |prefix for metrics (default "gen-statsd")                        |
  |RUN_TIME       |how long to run, 0=forever                                       |
  |SAMPLE_RATE    |sampling rate (default 0)                                        |
  |SPAWN_DRIFT    |spread new agent generation by 0-n seconds (default 10)          |
  |STATSD_HOSTS   |comma separated list of ip:port:proto for statsD host(s)         |
  |TAG_FORMAT     |format of the tags to send. accepted values "datadog" or "influx |
  |TAGS           |list of K:V comma separated tags. Example: key1:tag1,key2:tag2   |
  |TIMERS         |number of timers for each agent to hold (default 20)             |
  |VALUE_MAX      |maximum value to send (default 100)                              |
  |VALUE_MIN      |minimum value to send (default 0)                                |
  |VERSION        |show version information                                         |

## Releases

### Binary

[Releases](https://github.com/circonus-labs/gen-statsd/releases)

### Docker

See [Docker Hub](https://hub.docker.com/r/circonus/gen-statsd) for additional details.

* Be sure to have Docker installed on your machine
* Run: 
  * `docker pull circonus/gen-statsd:latest`
  * `docker run -e STATSD_HOST=127.0.0.1:8125 -e AGENTS=1 ... circonus/gen-statsd:latest`
* Note: You can also use the docker flag `--env-file` and provide the path to a file with the environment variables from above defined.

## Installation

### From Source

- Be sure to have `GOBIN` set
- Be sure that the location of GOBIN is in your `PATH`

1. Clone this repo
1. Navigate to the repo in your filesystem
1. Run `go install`

## Contributing

[Contributing Guide](./CONTRIBUTING.md)
