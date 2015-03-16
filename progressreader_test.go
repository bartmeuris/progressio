package progressio

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
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

func printProgress(msg string, t *testing.T, ch <-chan Progress) {
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

func TestWriter(t *testing.T) {
	r := getReader()
	w, ch := NewProgressWriter(getWriter(), -1)

	go printProgress("TestWriter", t, ch)

	io.Copy(w, r)
	t.Logf("Copy done\n")
}

func TestWriterSize(t *testing.T) {
	r := getReader()
	w, ch := NewProgressWriter(getWriter(), -1)

	go printProgress("TestWriterSize", t, ch)

	io.Copy(w, r)
	t.Logf("Copy done\n")
}

func TestReader(t *testing.T) {
	r, ch := NewProgressReader(getReader(), -1)
	w := getWriter()

	go printProgress("TestReader", t, ch)

	io.Copy(w, r)
	t.Logf("Copy done\n")
}

func TestReaderSize(t *testing.T) {
	r, ch := NewProgressReader(getReader(), bufSize)
	w := getWriter()

	go printProgress("TestReaderSize", t, ch)

	io.Copy(w, r)
	t.Logf("Copy done\n")
}
