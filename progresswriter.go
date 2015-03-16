package progressio

import "io"

// Copy functionality of ioutil.NopCloser, but for Writers
type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }
func getNopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

// ProgressWriter is a struct representing an io.WriterCloser, which sends back progress
// feedback over a channel
type ProgressWriter struct {
	w io.WriteCloser
	ioProgress
}

// NewProgressWriter creates a new ProgressWriter object based on the io.Writer and the
// size you specified. Specify a size <= 0 if you don't know the size.
func NewProgressWriter(w io.Writer, size int64) (*ProgressWriter, <-chan Progress) {
	if w == nil {
		return nil, nil
	}
	wc, ok := w.(io.WriteCloser)
	if !ok {
		wc = getNopWriteCloser(w)
	}
	ret := &ProgressWriter{wc, *mkIoProgress(size)}
	return ret, ret.ch
}

// Write wraps the io.Writer Write function to also update the progress.
func (p *ProgressWriter) Write(b []byte) (n int, err error) {
	n, err = p.w.Write(b[0:])
	p.updateProgress(int64(n))
	return
}

// Close wraps the io.WriterCloser Close function to clean up everything. ProgressWriter
// objects should always be closed to make sure everything is cleaned up.
func (p *ProgressWriter) Close() (err error) {
	err = p.w.Close()
	p.stopProgress()
	return
}
