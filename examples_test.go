package progressio

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

const throttleTime = 15 * time.Millisecond
const bufSize = 15000000

type throttleWriter struct {
	lastTime time.Time
}

func (t *throttleWriter) Write(b []byte) (n int, err error) {
	sleep := time.Since(t.lastTime)
	if t.lastTime.IsZero() {
		sleep = 0
	}
	if sleep < throttleTime {
		<-time.After(throttleTime - sleep)
	}
	t.lastTime = time.Now()
	return len(b), nil
}

type throttleReader struct {
	i        io.Reader
	lastTime time.Time
}

func (t *throttleReader) Read(b []byte) (n int, err error) {
	sleep := time.Since(t.lastTime)
	if t.lastTime.IsZero() {
		sleep = 0
	}
	if sleep < throttleTime {
		<-time.After(throttleTime - sleep)
	}
	t.lastTime = time.Now()
	return t.i.Read(b)
}

func getWriter() io.Writer {
	return &throttleWriter{}
}
func getReader() io.Reader {
	ibuf := make([]byte, bufSize)
	return &throttleReader{
		i: bytes.NewBuffer(ibuf),
	}
}

func printProgress(msg string, ch <-chan Progress) {
	cs := ""
	p := Progress{}
	for p = range ch {
		ps := msg + ": " + p.String()

		if len(cs) < len(ps) {
			cs = strings.Repeat(" ", len(ps))
		}
		fmt.Printf("\r%s\r%s", cs, ps)
	}
	fmt.Printf("\n%s\n", p.String())
}

// ExampleWriter is an example of using the progressio package with an
// io.Writer without knowing the amount of bytes to be processed.
func ExampleWriter() {
	r := getReader()
	w, ch := NewProgressWriter(getWriter(), -1)

	go printProgress("TestWriter", ch)

	io.Copy(w, r)
	fmt.Printf("Copy done\n")
}

// ExampleWriterSize is an example of using the progressio package with an
// io.Writer while knowing the expected amount of bytes to be processed.
func ExampleWriterSize() {
	r := getReader()
	w, ch := NewProgressWriter(getWriter(), -1)

	go printProgress("TestWriterSize", ch)

	io.Copy(w, r)
	fmt.Printf("Copy done\n")
}

// ExampleReader is an example of using the progressio package with an
// io.Reader without knowing the amount of bytes to be processed.
func ExampleReader() {
	r, ch := NewProgressReader(getReader(), -1)
	w := getWriter()

	go printProgress("TestReader", ch)

	io.Copy(w, r)
	fmt.Printf("Copy done\n")
}

// ExampleReaderSize is an example of using the progressio package with an
// io.Reader while knowing the expected amount of bytes to be processed.
func ExampleReaderSize() {
	r, ch := NewProgressReader(getReader(), bufSize)
	w := getWriter()

	go printProgress("TestReaderSize", ch)

	io.Copy(w, r)
	fmt.Printf("Copy done\n")
}
