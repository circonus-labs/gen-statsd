# gen-statsd

A StatsD Traffic Generator

## Getting Started

### Prerequisities

In order to run this container you'll need docker installed.

* [Windows](https://docs.docker.com/windows/started)
* [OS X](https://docs.docker.com/mac/started/)
* [Linux](https://docs.docker.com/linux/started/)

### Usage

#### Container Parameters

Running with command-line parameters:

```shell
docker run circonuslabs/gen-statsd:latest -agents 1 -statsd-host <statsd server IP>:8125 -prefix test -counters 1 -gauges 1 -timers 1 -protocol udp -spawn-drift 10 -tag-format datadog -tags key1:value1,key2:value2 
```

Printing the version:

```shell
docker run circonuslabs/gen-statsd:latest -version
```

Running in detached mode:

```shell
docker run -d circonuslabs/gen-statsd:latest -agents 1 ...
```

#### Environment Variables

* `AGENTS` - max number of agents to run concurrently (default 10)
* `COUNTERS` - number of counters for each agent to hold (default 50)
* `FLUSH_INTERVAL` - how often to flush metrics (default 10s)
* `GAUGES` - number of gauges for each agent to hold (default 30)
* `PREFIX` - prefix for metrics (default "gen-statsd")
* `PROTOCOL` - network protocol to use, tcp or udp (default "udp")
* `SPAWN_DRIFT` - spread new agent generation by 0-n seconds (default 10)
* `STATSD_HOST` - address of statsD host (default "localhost:8125")
* `TAG_FORMAT` - format of the tags to send. accepted values "datadog" or "influx
* `TAGS` - list of K:V comma separated tags. Example: key1:tag1,key2:tag2
* `TIMERS` - number of timers for each agent to hold (default 20)
* `VERSION` - show version information

### Image Variants

#### x86_64

[Distroless](https://github.com/GoogleContainerTools/distroless) Linux based x86_64 image.

* `circonuslabs/gen-statsd:latest`
* `circonuslabs/gen-statsd:1.0`
* `circonuslabs/gen-statsd:1.0.0`

#### ARM64

Alpine ARM64v8 based image.

* `circonuslabs/gen-statsd:latest-arm64`
* `circonuslabs/gen-statsd:1.0-arm64`
* `circonuslabs/gen-statsd:1.0.0-arm64`

## Find Us

* [GitHub](https://github.com/circonuslabs/gen-statsd)

## Contributing

Please read [CONTRIBUTING](https://github.com/circonus-labs/gen-statsd/blob/main/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the 
[tags on this repository](https://github.com/circonuslabs/gen-statsd/tags). 

## Authors

* **Justin DeFrank** - [Circonus](https://github.com/circonuslabs)

See also the list of [contributors](https://github.com/circonuslabs/gen-statsd/contributors) who 
participated in this project.

## License

This project is licensed under the BSD-3 License - see the [LICENSE](https://github.com/circonus-labs/gen-statsd/blob/main/LICENSE) file for details.