package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/omen-/httplog"

	"github.com/omen-/httplog/commonformat"
	"github.com/omen-/httplog/logfile"
)

const (
	minTimeBetweenRequests = 200
	maxTimeBetweenRequests = 1000
)

var (
	commonStatusCodes = []int{200, 301, 302, 404, 418, 500, 502}
	commonMethods     = []string{"GET", "POST", "PUT", "DELETE"}
	websiteResources  = []string{"/toto", "/toto/titi", "/foo", "/bar", "/"}
)

func main() {
	outputFile := flag.String("out", "access.log", "Output file path")
	flag.Parse()

	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)

	for {
		logWriter := logfile.NewWriter(*outputFile, commonformat.LogSerializer{})
		logWriter.WriteLogEntry(randomLogEntry(randomGenerator))
		timeUntilNextRequest := time.Duration(randomGenerator.Intn(maxTimeBetweenRequests-minTimeBetweenRequests) + minTimeBetweenRequests)
		time.Sleep(timeUntilNextRequest * time.Millisecond)
	}
}

func randomLogEntry(randomGenerator *rand.Rand) httplog.LogEntry {
	var logEntry httplog.LogEntry

	logEntry.IP = randomIPV4(randomGenerator)
	logEntry.Identity = "-"
	logEntry.UserID = "-"
	logEntry.DateTime = time.Now().Format(commonformat.TimeLayout)
	logEntry.Request = randomRequest(randomGenerator)
	logEntry.StatusCode = randomStatusCode(randomGenerator)
	logEntry.BytesSent = randomGenerator.Int63n(4096)

	return logEntry
}

func randomIPV4(randomGenerator *rand.Rand) string {
	return fmt.Sprintf("%v.%v.%v.%v", randomGenerator.Intn(256), randomGenerator.Intn(256), randomGenerator.Intn(256), randomGenerator.Intn(256))
}

func randomRequest(randomGenerator *rand.Rand) string {
	return fmt.Sprintf("%v %v HTTP/1.%v", commonMethods[randomGenerator.Intn(len(commonMethods))],
		websiteResources[randomGenerator.Intn(len(websiteResources))], randomGenerator.Intn(2))
}

func randomStatusCode(randomGenerator *rand.Rand) int {
	return commonStatusCodes[randomGenerator.Intn(len(commonStatusCodes))]
}
