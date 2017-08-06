package tracer

import (
	"log"
	"os"

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

func Get() *statsd.Client {
	return tracer
}
