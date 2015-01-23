package main

import (
	"fmt"
	"net/http"
	"time"
)

func Collector(w http.ResponseWriter, r *http.Request) {

	// Parse the delay.
	delay, err := time.ParseDuration("1s")
	if err != nil {
		http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		// Now, we take the delay, and the person's name, and make a WorkRequest out of them.
		work := WorkRequest{Name: "max", Delay: delay}
		// Push the work onto the queue.
		Dispatch.WorkQueue <- work
		fmt.Println("Work request queued")

		return
	}()

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return
}

func DispatchReload(w http.ResponseWriter, r *http.Request) {

	Dispatch.Stop()

	err := NewDispatcher(50)
	if err != nil {
		fmt.Println(err.Error())
	}

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return
}
