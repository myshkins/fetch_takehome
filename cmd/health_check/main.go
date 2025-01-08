package main

import (
	"flag"
)

import (
	"github.com/myshkins/fetch_takehome/internal/health_check_client"
)

func main() {
	endpointConfigFilePath := flag.String("config-file", "../../input.yaml", "file path to the endpoint config yaml to be used")
  timeInterval := flag.Int("interval", 15, "time interval in seconds to check endpoints")
	flag.Parse()
  hc := health_check_client.NewHealthCheckClient(*endpointConfigFilePath, *timeInterval)
  hc.PingEndpoints()
}
