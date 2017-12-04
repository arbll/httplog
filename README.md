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
    	An alert will be generated if the trafic for the past <aperiod> minutes exceed the given threshold. (default 2m0s)
  -logpath string
    	Path to the common log format file. (default "access.log")
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

## Dependencies

**[hpcloud/tail](https://github.com/hpcloud/tail)**

> Go package for reading from continously updated files (tail -f) 

Reading continously updated files is platform dependent and a hard task in general. This library provide a cross-platform easy to use way of doing so. It is also the most widely used for this task.

---

**[gizak/termui](https://github.com/gizak/termui)**

> Golang terminal dashboard 

Pretty terminal ui.

---

