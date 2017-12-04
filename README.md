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

```
> monitorhttplog -h

  -h	Show usage
  -aperiod duration
    	The traffic will be monitored [aperiod] time back to raise threshold alerts. (default 2m0s)
  -logpath string
    	Path to the common format log file. (default "access.log")
  -rperiod duration
    	Frequency at which reports will be generated. (default 10s)
  -threshold int
    	Traffic threshold after which an alert will be generated. (default 200)
```

### Generator

```
> generatehttplog -h
 
  -h    Show usage
  -out string
        Output file path (default "access.log")
```
