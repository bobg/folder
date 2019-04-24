package maildir

// TODO: add tests

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
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
	return &Maildir{name, append(files, files2...)}, nil
}

func getfiles(dir, subdir string) ([]string, error) {
	infos, err := ioutil.ReadDir(path.Join(dir, subdir))
	if err != nil {
		return nil, err
	}
	var result []string
	for _, info := range infos {
		result = append(result, path.Join(subdir, info.Name()))
	}
	return result, nil
}

func (m *Maildir) Message() (io.ReadCloser, error) {
	if len(m.files) == 0 {
		return nil, nil
	}
	file := path.Join(m.name, m.files[0])
	m.files = m.files[1:]
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "opening %s", file)
	}
	return f, nil
}
