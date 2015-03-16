# progressio
Library to get progress feedback from io.Reader and io.Writer objects

To install, run `go get github.com/bartmeuris/progressio`

# Useage

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
