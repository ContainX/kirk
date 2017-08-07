package tracer

import (
	"log"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

var tracer *statsd.Client
var initErr error

func init() {
	tracer, initErr = statsd.New(os.Getenv("HOST") + ":8125")

	tracer.Namespace = "kirk."

	if initErr != nil {
		log.Fatal("Error initializing statsd agent:", initErr)
	} else {
		log.Println("Datadog tracer initialized")
	}

}

// Timer, utiliy for measuring timing
type timerInstance struct {
	startTime time.Time
	metric    string
	tags      []string
}

func Timer(metric string, tags []string) timerInstance {
	t := timerInstance{
		startTime: time.Now(),
		metric:    metric,
		tags:      tags,
	}

	return t
}

func (t *timerInstance) End() {
	durationT := time.Since(t.startTime)
	ns := durationT.Nanoseconds()
	ms := ns / 1000000

	tracer.Histogram(t.metric, float64(ms), t.tags, 1)
}

// Return the tracer instance
func Get() *statsd.Client {
	return tracer
}
