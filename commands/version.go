package commands

import (
	"fmt"
	"time"
)

var lastUpdated time.Time

func init() {
	lastUpdated = time.Now()
}

func versionCommand() string {
	return fmt.Sprintf("I was last updated at: %v", lastUpdated)
}
