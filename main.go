package main

import (
	"net/http"
	"time"
)

func makePoop(w http.ResponseWriter, req *http.Request) {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "cannot make poop", http.StatusInternalServerError)
		return
	}
	for {
		select {
		case <-req.Context().Done():
			return
		case <-t.C:
			if _, err := w.Write([]byte("💩")); err != nil {
				return
			}
			f.Flush()
		}
	}
}

func main() {
	http.HandleFunc("/", makePoop)
	http.ListenAndServe(":8080", nil)
}
