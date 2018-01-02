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
	files, err := getfiles(name, "cur")
	if err != nil {
		return nil, err
	}
	files2, err := getfiles(name, "new")
	if err != nil {
		return nil, err
	}
	return &Maildir{name, append(files, files2...)}
}

func getfiles(dir, subdir string) ([]string, error) {
	infos, err := ioutil.ReadDir(path.Join(dir, subdir))
	if err != nil {
		return nil, err
	}
	var result []string
	for _, info := range infos {
		files = append(files, path.Join(subdir, info.Name()))
	}
	return files, nil
}

func (m *Maildir) Message() (io.Reader, func() error, error) {
	if len(m.files) == 0 {
		return nil, nil, nil
	}
	file := path.Join(m.name, m.files[0])
	m.files = m.files[1:]
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, fmt.Errorf("opening %s: %s", file, err)
	}
	close := func() error { return f.Close() }
	return f, close, nil
}
