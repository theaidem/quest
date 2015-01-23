package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/rakyll/gometry/http"
)

var (
	NWorkers = flag.Int("n", 5, "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
)

// clear && go build -o que ./*go && ./que -n=5
// for i in {1..1000}; do curl "localhost:8000/work"; done
func main() {
	// Parse the command-line flags.
	flag.Parse()

	// Start the dispatcher.
	fmt.Println("Starting the dispatcher")

	err := NewDispatcher(*NWorkers)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Register our collector as an HTTP handler function.
	fmt.Println("Registering the collector")
	http.HandleFunc("/work", Collector)
	http.HandleFunc("/dispatch", DispatchReload)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}
}
