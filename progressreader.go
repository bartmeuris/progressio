package progressio

import (
	"io"
	"io/ioutil"
	"os"
)

// ProgressReader is a struct representing an io.ReaderCloser, which sends back progress
// feedback over a channel
type ProgressReader struct {
	r io.ReadCloser
	ioProgress
}

// NewProgressFileReader creates a new ProgressReader based on a file. It teturns a
// ProgressReader object and a channel on success, or an error on failure.
func NewProgressFileReader(file string) (*ProgressReader, <-chan Progress, error) {
	f, ferr := os.Open(file);
	if ferr != nil {
		return nil, nil, ferr
	}
	// Get the filesize by seeking to the end of the file, and back to offset 0
	fsize, err := f.Seek(0, os.SEEK_END)
	if err != nil {
		return nil, nil, err
	}
	if _, err := f.Seek(0, os.SEEK_SET); err != nil {
		return nil, nil, err
	}
	io, ch := NewProgressReader(f, fsize)
	return io, ch, nil
}

// NewProgressReader creates a new ProgressReader object based on the io.Reader and the
// size you specified. Specify a size <= 0 if you don't know the size.
func NewProgressReader(r io.Reader, size int64) (*ProgressReader, <-chan Progress) {
	if r == nil {
		return nil, nil
	}
	rc, ok := r.(io.ReadCloser)
	if !ok {
		rc = ioutil.NopCloser(r)
	}
	ret := &ProgressReader{rc, *mkIoProgress(size)}
	return ret, ret.ch
}

// Read wraps the io.Reader Read function to also update the progress.
func (p *ProgressReader) Read(b []byte) (n int, err error) {
	n, err = p.r.Read(b)
	p.updateProgress(int64(n))
	return
}

// Close wraps the io.ReaderCloser Close function to clean up everything. ProgressReader
// objects should always be closed to make sure everything is cleaned up.
func (p *ProgressReader) Close() (err error) {
	err = p.r.Close()
	p.stopProgress()
	return
}
