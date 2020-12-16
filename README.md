# gen-statsd

[![Go Report Card](https://goreportcard.com/badge/github.com/circonus-labs/gen-statsd)](https://goreportcard.com/report/github.com/circonus-labs/gen-statsd)

gen-statsd is an experimental tool written in Go for load testing with StatsD metrics.

## Configuration

### Command-line Flags

```
  Usage of gen-statsd:
  -agents int
        max number of agents to run concurrently (default 10)
  -counters int
        number of counters for each agent to hold (default 50)
  -flush-interval duration
        how often to flush metrics (default 10s)
  -gauges int
        number of gauges for each agent to hold (default 30)
  -prefix string
        prefix for metrics (default "gen-statsd")
  -protocol string
        network protocol to use, tcp or udp (default "udp")
  -spawn-drift int
        spread new agent generation by 0-n seconds (default 10)
  -statsd-host string
        address of statsD host (default "localhost:8125")
  -tag-format string
        format of the tags to send. accepted values "datadog" or "influx"
  -tags string
        list of K:V comma separated tags. Example: key1:tag1,key2:tag2
  -timers int
        number of timers for each agent to hold (default 20)
  -version
        show version information
```

### Environment Variables

  | Variable      |                       Description                               |
  |:-------------:|:----------------------------------------------------------------|
  |AGENTS         |max number of agents to run concurrently (default 10)            |
  |COUNTERS       |number of counters for each agent to hold (default 50)           |
  |FLUSH_INTERVAL |how often to flush metrics (default 10s)                         |
  |GAUGES         |number of gauges for each agent to hold (default 30)             |
  |PREFIX         |prefix for metrics (default "gen-statsd")                        |
  |PROTOCOL       |network protocol to use, tcp or udp (default "udp")              |
  |SPAWN_DRIFT    |spread new agent generation by 0-n seconds (default 10)          |
  |STATSD_HOST    |address of statsD host (default "localhost:8125")                |
  |TAG_FORMAT     |format of the tags to send. accepted values "datadog" or "influx |
  |TAGS           |list of K:V comma separated tags. Example: key1:tag1,key2:tag2   |
  |TIMERS         |number of timers for each agent to hold (default 20)             |
  |VERSION        |show version information                                         |

## Releases

### Binary

[Releases](https://github.com/circonus-labs/gen-statsd/releases)

### Docker

See [Docker Hub](hub.docker.com/circonus/gen-statsd) for additional details.

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
