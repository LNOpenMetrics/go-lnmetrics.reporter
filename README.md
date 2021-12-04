<div align="center">
  <h1> :chart_with_upwards_trend: go-lnmetrics-reporter :bar_chart: </h1>

  <img src="https://github.com/OpenLNMetrics/lnmetrics.icons/blob/main/current/res/mipmap-xxxhdpi/ic_launcher.png" />

  <p>
    <strong> Reference implementation written in Go lang to collect and report of the lightning node metrics </strong>
  </p>

  <p>
   <a href="https://github.com/LNOpenMetrics/go-lnmetrics.reporter/actions">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/LNOpenMetrics/go-lnmetrics.reporter/Build%20and%20test%20Go?style=flat-square"/>
   </a>
   <a href="https://discord.gg/vFX989za">
    <img alt="Discord" src="https://img.shields.io/discord/913794833498394634?style=flat-square">
   </a>
   <a href="https://github.com/LNOpenMetrics/go-lnmetrics.reporter/releases">
    <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/LNOpenMetrics/go-lnmetrics.reporter?style=flat-square"/>
   </a>
  </p>
</div>

## Table of Content

- Introduction
- Install Procedure
- How to Use
- How to Contribute
- Build With
- License

## Introduction

go-lnmetrics.reporter is the reference implementation of a lnmetrics client (see definition of lnmetrics client [here](https://github.com/LNOpenMetrics/lnmetrics.rfc#terminology)). It is written in Go lang to use the power of the statically compiled binary that a user can  easily run on the target machine. This help to speedup the 
data collection process very easily.

In addition, this plugin support the following metrics, and you can also see the in the spec what data of your ln node are shared with the server.

- [Metric One](https://github.com/LNOpenMetrics/lnmetrics.rfc/blob/main/metrics/metric_1.md): See the definition on the specification [lnmetrics.spec](https://github.com/LNOpenMetrics/lnmetrics.rfc)

## Install Procedure

It is suggested to work with the tagged version that you can find in the [release page](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/releases).

The installation process required just to download the right binary from your host machine, and configure c-lightning (that it is the only implementation supported right now) to run the plugin and pass the server link to lightnind.

The configuration suggested is to use the config file, with the following content

```
... other stuff...
 
plugin=/path/binary/go-lnmetrics-{archirecture}
lnmetrics-urls=https://api.lnmetrics.info/query
```

Where the `lnmetrics-urls` give the possibility to specify the server where the plugin need to report the data if any. You can specify more that one server with just
append another url divided by a comma.

## How to Use

After running the plugin you will have the possibility to run the following rpc command from lightning-cli, and them are described below:

- `metric_one start end`: RPC command that give you the possibility to query the internal db and make access to the metric data collected by the plugin. The start and end need to be a string that is the timestamp or you can use the following query to make some particular query:
  - `metric_one start="now"`: Give you the possibility to query the metric data that the plugin have in memory;
  - `metric_one start="last"`: Give you the possibility to query the metric data that the plugin committed to the server last time.
- `lnmetrics-info`: RPC command that give you access to the plugin information, like version, go version and architecture this will be useful when there is some bug
report or just consult the version of the plugin that the user is running.

## How to Contribute

You can contribute in three different way, and them are described below:

- Bug and feature request through Github discussion/issue or discord channel;
- Bug fixing through PR;
- New metric support, this required to follow the [lnmetrics.spec guide line](https://github.com/LNOpenMetrics/lnmetrics.rfc#how-propose-a-new-metric)

In addition, if you want build the project or you can start to play with it, you need the golang compiler (suggested the last one) and the golangci (see Build With section)
to compile the code with the make command.

## Build With
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- [glightning](https://github.com/vincenzopalazzo/glightning)
- [golangci](https://golangci-lint.run/)

## License

TODO
