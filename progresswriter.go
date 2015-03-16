package progressio

import "io"
// Copy functionality of ioutil.NopCloser, but for Writers
type nopWriteCloser struct { io.Writer }
func (nopWriteCloser) Close() error { return nil }
func getNopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

type ProgressWriter struct {
	w io.WriteCloser
	ioProgress
}

func NewProgressWriter(w io.Writer, size int64) (*ProgressWriter, <- chan Progress) {
	if w == nil {
		return nil, nil
	}
	wc, ok := w.(io.WriteCloser)
	if !ok {
		wc = getNopWriteCloser(w)
	}
	ret := &ProgressWriter{wc, *mkIoProgress( size ) }
	return ret, ret.ch
}

func (p *ProgressWriter) Write(b []byte) (n int, err error) {
	n, err = p.w.Write(b[0:])
	p.updateProgress(int64(n))
	return
}

func (p *ProgressWriter) Close() (err error) {
	err = p.w.Close()
	p.stopProgress()
	return
}

