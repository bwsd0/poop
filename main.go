// poop sends an endless stream of HTTP chunk transfer encoded poop.
//
// Usage:
//
//	poop [-a addr] [-d delay]
//
// BUG: poop cannot be sent using HTTP2 as the specification forbids use
// chunked-transfer encoding. See: RFC 7540, 8.1.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	addr = flag.String("addr", ":8080", "HTTP listen address")
)

func makePoop(w http.ResponseWriter, req *http.Request) {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "cannot make poop", http.StatusInternalServerError)
		return
	}

	// Browsers may buffer chunks instead of rendering the data as it arrives; To
	// mitigate this, write a partial response with a small payload (1KB) of
	// non-renderable Unicode and prevent MIME-type sniffing which may also cause
	// buffering.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Zero Width Space ("\u200B")
	zws := []byte{0xE2, 0x80, 0x8B}
	if _, err := w.Write(bytes.Repeat(zws, 1024/len(zws))); err != nil {
		panic(err)
	}

	var n int
	var err error
	for i := 0; i <= 1024; i += n {
		select {
		case <-req.Context().Done():
			return
		case <-t.C:
			if n, err = w.Write([]byte{0xF0, 0x9F, 0x92, 0xA9}); err != nil {
				return
			}
			f.Flush()
		}
	}
}

const usageLine = `usage: poop [-a addr]`

func main() {
	flag.Parse()
	if flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, usageLine)
		os.Exit(2)
	}

	http.HandleFunc("/", makePoop)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
