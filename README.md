# HTTP log monitor

[![GoDoc](https://godoc.org/github.com/Omen-/httplog?status.svg)](https://godoc.org/github.com/Omen-/httplog) [![Go Report Card](https://goreportcard.com/badge/github.com/Omen-/httplog)](https://goreportcard.com/report/github.com/Omen-/httplog)

![Screenshot](/_assets/screenshot.png)

## Installation

This project requires Go 1.8 or later.

```bash
go get -u github.com/omen-/httplog/cmd/generatehttplog
go get -u github.com/omen-/httplog/cmd/monitorhttplog
```

## Documentation

Documentation can be found on [godoc](https://godoc.org/github.com/Omen-/httplog).

## Usage

If you just want to try the program just create a directory and run in any order :

```
generatehttplog
```
```
monitorhttplog
```

### Monitor

```
> monitorhttplog -h

  -h    Show usage
  -aperiod duration
        An alert will be generated if the traffic for the past <aperiod> minutes exceed the given threshold. (default 2m0s)
  -logpath string
        Path to the common log format file. (default "access.log")
  -rperiod duration
        Frequency at which reports will be generated. (default 10s)
  -threshold int
        Traffic threshold after which an alert will be generated. (default 150)
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

> Go package for reading from continuously updated files (tail -f) 

Reading continuously updated files is platform dependent and a complex task in general. This library provides a cross-platform and easy to use way of doing so. It is also the most widely used for this task.

---

**[gizak/termui](https://github.com/gizak/termui)**

> Golang terminal dashboard 

Pretty terminal ui.

---

## Design decisions & Possible improvements

+ Right now, the alert check is invoked when a new log entry is consumed. I made this design decision to ensure that alerts are shown as triggered at the time the server exceeded the threshold. It seemed important to me to trust the time of the log more than our system time since there might be delays (server flushing log file periodically, monitor overloaded, ...) and that may cause "fake" alerts to trigger in certain cases. Concretely this leads to a scenario where a "back to normal" alert is not triggered if the traffic stop abruptly and nothing is logged for 2+ minutes.

+ This principle of only trusting the time from the log files is not applied on the 10s reports. It would be a big improvement and would allow loading past log data.

+ Support more log formats. I made sure to split monitoring and reading/parsing responsibilities to allow for future improvements without too much refactoring.

+ Monitor more data: Status codes, Requests by IP addresses, ...
