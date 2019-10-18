// poop sends an endless stream of HTTP chunk transfer encoded poop.
//
// Usage:
//
//  go get [-u] github.com/bwasd/poop
//  $BROWSER http://localhost:8080
//
// In a terminal emulator with Unicode support, with a font that
// includes the poop emoji glyph, invoke `curl(1)` with the `-N` option
// to disable buffering and enjoy this app from the comfort of your
// terminal.
package main

import (
	"bytes"
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

	// Browsers may buffer chunks instead of rendering the data as it arrives;
	// To mitigate this, write a partial response with a small payload (1KB) of
	// non-renderable Unicode and prevent MIME-type sniffing which may also
	// cause buffering.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Zero Width Space ("\u200B")
	zws := []byte{0xE2, 0x80, 0x8B}
	bufsz := 1024
	buf := make([]byte, bufsz)
	buf = bytes.Repeat(zws, bufsz/len(zws))
	if _, err := w.Write(buf); err != nil {
		panic(err)
	}

	for {
		// FIXME: returning an indeterminate-length data stream can be used as an
		// exploitable resource exhaustion vector.
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
