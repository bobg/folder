package maildir

// TODO: add tests

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type Maildir struct {
	name  string
	files []string
}

func New(name string) (*Maildir, error) {
	cur := path.Join(name, "cur")
	infos, err := ioutil.ReadDir(cur)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, info := range infos {
		files = append(files, info.Name())
	}
	return &Maildir{name, files}, nil
}

func (m *Maildir) Message() (io.Reader, func() error, error) {
	if len(m.files) == 0 {
		return nil, nil, nil
	}
	file := path.Join(m.name, "cur", m.files[0])
	m.files = m.files[1:]
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, fmt.Errorf("opening %s: %s", file, err)
	}
	close := func() error { return f.Close() }
	return f, close, nil
}
