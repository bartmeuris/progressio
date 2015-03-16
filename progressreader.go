package progressio

import (
	"os"
	"io"
	"io/ioutil"
)

type ProgressReader struct {
	r io.ReadCloser
	ioProgress
}

func NewProgressFileReader(file string) (*ProgressReader, <- chan Progress, error) {
	if f, ferr := os.Open(file); ferr == nil {
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
	} else {
		return nil, nil, ferr
	}
}

func NewProgressReader(r io.Reader, size int64) (*ProgressReader, <- chan Progress) {
	if r == nil {
		return nil, nil
	}
	rc, ok := r.(io.ReadCloser)
	if !ok {
		rc = ioutil.NopCloser(r)
	}
	ret := &ProgressReader{rc, *mkIoProgress( size ) }
	return ret, ret.ch
}

func (p *ProgressReader) Read(b []byte) (n int, err error) {
	n, err = p.r.Read(b)
	p.updateProgress(int64(n))
	return
}

func (p *ProgressReader) Close() (err error) {
	err = p.r.Close()
	p.stopProgress()
	return
}

