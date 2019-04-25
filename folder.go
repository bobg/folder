package folder

import (
	"bufio"
	"bytes"
	"io"

	"github.com/bobg/uncompress"

	"github.com/bobg/folder/v3/maildir"
	"github.com/bobg/folder/v3/mbox"
	"github.com/bobg/folder/v3/tar"
)

// Folder allows reading messages one by one.
type Folder interface {
	// Message returns an io.ReadCloser for the next message in the folder.
	// After the last message,
	// Message returns a nil ReadCloser.
	// After calling Message,
	// it is an error to call it again before first closing the ReadCloser.
	Message() (io.ReadCloser, error)

	// Close releases resources held by the folder.
	Close() error
}

// Open opens the given pathname as a mail folder if possible,
// trying one type after another and uncompressing on the fly if needed.
func Open(name string) (Folder, error) {
	m, err := maildir.New(name)
	if err != nil {
		return m, nil
	}

	r, err := uncompress.OpenFile(name)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			r.Close()
		}
	}()

	b := bufio.NewReader(r)
	from, err := b.Peek(5)
	if err != nil {
		return nil, err
	}

	rc := &readCloser{
		b: b,
		r: r,
	}

	if bytes.Equal(from, []byte("From ")) {
		return mbox.New(rc)
	}

	return tar.New(rc), nil
}

type readCloser struct {
	b *bufio.Reader
	r io.ReadCloser
}

func (rc *readCloser) Read(buf []byte) (int, error) {
	return rc.b.Read(buf)
}

func (rc *readCloser) Close() error {
	return rc.r.Close()
}
