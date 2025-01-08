package health_check_client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

type endpoint struct {
	Name    string            `yaml:"name"`
	Url     string            `yaml:"url"`
	Method  string            `yaml:"method,omitempty"`
	Body    string            `yaml:"body,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type HealthCheckClient struct {
  httpClient *http.Client
  endpoints []endpoint
  timeInterval int
  stats map[string]map[string]int  // eg. {"endpoint.Name":{"up":5}}
}

func NewHealthCheckClient(fp string, t int) *HealthCheckClient {
  var hc HealthCheckClient
  hc.httpClient = &http.Client{
    Timeout: 500 * time.Millisecond,
  }
  hc.endpoints = parseEndpointConfig(fp)
  hc.timeInterval = t
  hc.stats = make(map[string]map[string]int)

  // initialize inner map for each endpoint
  for _, e := range hc.endpoints {
    hc.stats[e.Name] = make(map[string]int)
  }
  return &hc
}

func parseEndpointConfig(fp string) []endpoint {
	endpoints := []endpoint{}

	data, err := os.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &endpoints)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return endpoints
}

func (hc *HealthCheckClient) printStats() {
  for _, v := range hc.endpoints {
    total := hc.stats[v.Name]["up"] + hc.stats[v.Name]["down"]
    a := (float64(hc.stats[v.Name]["up"]) / float64(total)) * float64(100)
    a = math.Round(a)
    fmt.Printf("%s has %+v%% availability percentage\n", v.Name, a)
  }
}

func (hc *HealthCheckClient) PingEndpoints() {
	ticker := time.NewTicker(time.Duration(hc.timeInterval) * time.Second)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("\nStarting endpoint healthcheck. Pinging every %d seconds. Press Ctrl+C to stop\n", hc.timeInterval)
  
  err := hc.ping()
  if err != nil {
    log.Printf("error pinging endpoints: ", err)
  }

	for {
		select {
		case <-ticker.C:
      err := hc.ping()
      if err != nil {
        log.Printf("error pinging endpoints: %v", err)
      }
		case <-sigChan:
			fmt.Println("\nReceived interrupt signal. Exiting...\n")
			return
		}
	}
}

func formRequestBody(endpoint endpoint) io.Reader {
  body := bytes.NewBuffer(nil)
  return body
}

func formRequest(endpoint endpoint) *http.Request {
  method := "GET"
  body := formRequestBody(endpoint)
  if endpoint.Method != "" {
    method = endpoint.Method
  }
  req, err := http.NewRequest(method, endpoint.Url, body)
  if err != nil {
    log.Printf("error forming request: %v", err)
  }

  if endpoint.Headers != nil {
    for k,v := range endpoint.Headers {
      req.Header.Add(k, v)
    }
  }
  return req
}

func (hc *HealthCheckClient) ping() error {
  fmt.Println("...")
	for _, endpoint := range hc.endpoints {
    req := formRequest(endpoint)
    
    resp, err := hc.httpClient.Do(req)
		if err != nil {
      hc.stats[endpoint.Name]["down"]++
      continue
		}
		defer resp.Body.Close()

    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
      hc.stats[endpoint.Name]["up"]++
    } else {
      hc.stats[endpoint.Name]["down"]++
    }
	}
  hc.printStats()
  return nil
}
