package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/omen-/httplog/pkg/commonformat"
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
	help := flag.Bool("h", false, "Show usage")
	outputFile := flag.String("out", "access.log", "Output file path")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)

	var step = math.Pi / 2 / 100
	x := .0

	for {
		logWriter := commonformat.NewWriter(*outputFile, &commonformat.LogSerializer{})
		logWriter.WriteLogEntry(randomLogEntry(randomGenerator))
		x += step
		timeUntilNextRequest := math.Abs(math.Sin(x)*(maxTimeBetweenRequests-minTimeBetweenRequests)) + minTimeBetweenRequests
		time.Sleep(time.Duration(timeUntilNextRequest) * time.Millisecond)
	}
}

func randomLogEntry(randomGenerator *rand.Rand) commonformat.LogEntry {
	var logEntry commonformat.LogEntry

	logEntry.IP = randomIPV4(randomGenerator)
	logEntry.Identity = "-"
	logEntry.UserID = "-"
	logEntry.Time = time.Now()
	logEntry.Request = randomRequest(randomGenerator)
	logEntry.StatusCode = randomStatusCode(randomGenerator)
	logEntry.BytesSent = randomGenerator.Int63n(4096)

	return logEntry
}

func randomIPV4(randomGenerator *rand.Rand) string {
	return fmt.Sprintf("%v.%v.%v.%v", randomGenerator.Intn(256), randomGenerator.Intn(256), randomGenerator.Intn(256), randomGenerator.Intn(256))
}

func randomRequest(randomGenerator *rand.Rand) commonformat.Request {
	return commonformat.Request{
		Method:      commonMethods[randomGenerator.Intn(len(commonMethods))],
		Resource:    websiteResources[randomGenerator.Intn(len(websiteResources))],
		HTTPVersion: "HTTP/1.1",
	}
}

func randomStatusCode(randomGenerator *rand.Rand) int {
	return commonStatusCodes[randomGenerator.Intn(len(commonStatusCodes))]
}
