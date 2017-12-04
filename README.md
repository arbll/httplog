# HTTP log monitor

[![GoDoc](https://godoc.org/github.com/Omen-/httplog?status.svg)](https://godoc.org/github.com/Omen-/httplog) [![Go Report Card](https://goreportcard.com/badge/github.com/Omen-/httplog)](https://goreportcard.com/report/github.com/Omen-/httplog)

## Installation

This project requires Go 1.8 or later.

```bash
go get -u github.com/omen-/httplog/cmd/generatehttplog
go get -u github.com/omen-/httplog/cmd/monitorhttplog
```

## Documentation

Documentation can be found on [godoc](https://godoc.org/github.com/Omen-/httplog).

## Usage

### Monitor

```bash
> monitorhttplog -h

  -h	Show usage
  -aperiod duration
    	The traffic will be monitored [aperiod] time beck to raise threshold alerts. (default 2m0s)
  -logpath string
    	Path to the common format log file. (default "access.log")
  -rperiod duration
    	Frequency at which reports will be generated. (default 10s)
  -threshold int
    	Traffic threshold after which an alert will be generated. (default 200)
```
go install github.com/omen-/httplog/cmd/generatehttplog
go install github.com/omen-/httplog/cmd/monitorhttplog