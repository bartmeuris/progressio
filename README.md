# progressio

[![GoDoc](https://godoc.org/github.com/bartmeuris/progressio?status.svg)](http://godoc.org/github.com/bartmeuris/progressio)

Library to get progress feedback from io.Reader and io.Writer objects

To install, run `go get github.com/bartmeuris/progressio`

## About

This package was created since most progress packages out-there seem to
directly want to output stuff to a console mixing UI and work locig, 
work in a non-Go way (callback functions), or can only be used for 
specific scenario's like file downloads from the web.

This package just wraps standard io.Reader and io.Writer objects so anything
that uses these can give you progress feedback. It attempts to do all the
heavy lifting for you:

* updates are throttled to 1 per 100ms (10 per second)
* Precalculates things (if possible) like:
  * Speed in bytes/sec of the last few operations
  * Average speed in bytes/sec since the start of the operation
  * Remaining time
  * Percentage

Some of these statistics are not available if the size was not specified up front.

## Examples

```
import (
  "io"
  "github.com/bartmeuris/progressio"
)

// io.Copy wrapper to specify the size and show copy progress.
func copyProgress(w io.Writer, r io.Reader, size int64) (written int64, err error) {
  pw, ch := progressio.NewProgressWriter(w, size)
  defer pw.Close()
  
  // Launch a Go-Routine reading from the progress channel
  go func() {
    for p := range ch {
      fmt.Printf("\rProgress: %s", p.String())
    }
    fmt.Printf("\nDone\n")
  }
  
  // Copy the data from the reader to the new writer
  return io.Copy(pw, r)
}
```

## TODO

* Add tests
* Clean up documentation
* Document byte/duration formatters
* Extract byte/duration formatters and place in separate library (?)

