package tar

import (
	atar "archive/tar"
	"io"
	"io/ioutil"
)

type Tar struct {
	t *atar.Reader
	r io.ReadCloser
}

func New(r io.ReadCloser) *Tar {
	return &Tar{
		t: atar.NewReader(r),
		r: r,
	}
}

func (t *Tar) Message() (io.ReadCloser, error) {
	for {
		h, err := t.t.Next()
		if err == io.EOF {
			return nil, nil
		}

		if h.Typeflag == atar.TypeReg {
			return ioutil.NopCloser(t.t), nil
		}
	}
}

func (t *Tar) Close() error {
	return t.r.Close()
}
