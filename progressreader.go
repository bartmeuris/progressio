package progressio

import "io"
import "io/ioutil"

type ProgressReader struct {
	r io.ReadCloser
	ioProgress
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
	err = p.Close()
	p.stopProgress()
	return
}

