package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/theaidem/quest"
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

	err := quest.NewDispatcher(*NWorkers)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Register our collector as an HTTP handler function.
	fmt.Println("Registering the collector")
	http.HandleFunc("/work", Collector)
	http.HandleFunc("/dispatch", DispatchReload)

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}

}

func Collector(w http.ResponseWriter, r *http.Request) {

	// Parse the delay.
	delay, err := time.ParseDuration("1s")
	if err != nil {
		http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		// Now, we take the delay, and the person's name, and make a WorkRequest out of them.
		work := quest.WorkRequest{Name: "max", Delay: delay}
		// Push the work onto the queue.
		quest.Dispatch.WorkQueue <- work
		fmt.Println("Work request queued")

		return
	}()

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return
}

func DispatchReload(w http.ResponseWriter, r *http.Request) {

	quest.Dispatch.Stop()

	err := quest.NewDispatcher(50)
	if err != nil {
		fmt.Println(err.Error())
	}

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return
}
