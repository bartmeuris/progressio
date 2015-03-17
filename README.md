# progressio

[![GoDoc](https://godoc.org/github.com/bartmeuris/progressio?status.svg)](http://godoc.org/github.com/bartmeuris/progressio)

Library to get progress feedback from io.Reader and io.Writer objects

To install, run `go get github.com/bartmeuris/progressio`

## About

This package was created since most progress packages out-there seem to
directly want to output stuff to a console mixing UI and work logic, 
work in a non-Go way (callback functions), or can only be used for 
specific scenario's like file downloads from the web.

This package provides wrappers around standard io.Reader and io.Writer objects
which send back a progressio.Progress struct over a channel, so anything
that uses standard io.Reader/io.Writer objects can give you progress feedback.
It attempts to do all the heavy lifting for you:

* updates are throttled to 1 per 100ms (10 per second)
* Precalculates things (if possible) like:
  * Speed in bytes/sec of the last few operations
  * Average speed in bytes/sec since the start of the operation
  * Remaining time
  * Percentage

Some of these statistics are not available if the size was not specified up front.

## Progress object

### Layout

```
type Progress struct {
▸   Transferred int64         // Transferred data in bytes
▸   TotalSize   int64         // Total size of the transfer in bytes. <= 0 if size is unknown.
▸   Percent     float64       // If the size is known, the progress of the transfer in %
▸   SpeedAvg    int64         // Bytes/sec average over the entire transfer
▸   Speed       int64         // Bytes/sec of the last few reads/writes
▸   Remaining   time.Duration // Estimated time remaining, only available if the size is known.
▸   StartTime   time.Time     // When the transfer was started
▸   StopTime    time.Time     // only specified when the transfer is completed: when the transfer was stopped
}

```

### Functions

The progressio.Progress object has at the moment only one function, the
String() function to return the `string` representation of the object.

## Example

```
import (
  "io"
  "github.com/bartmeuris/progressio"
)

// io.Copy wrapper to specify the size and show copy progress.
func copyProgress(w io.Writer, r io.Reader, size int64) (written int64, err error) {
  
  // Wrap your io.Writer:
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

