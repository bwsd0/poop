poop sends an endless stream of HTTP chunk transfer encoded poop.

Its main use is a "hello world" sanity check or post-install test when
updating versions of Go or Docker, for example.

Usage:
```sh
go get [-u] github.com/bwasd/poop
poop &
$BROWSER http://localhost:8080
```

In a terminal emulator that supports Unicode with a font that includes
the poop emoji glyph, invoke `curl(1)` with the `-N` option to disable
buffering.
