// Package mbox implements parsing of mbox-style mail folders.

package mbox

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/bobg/chanrw"
)

// Mbox is a parser for a Unix mbox-style mail folder.
type Mbox struct {
	scanner *bufio.Scanner
	eof     bool
}

func New(r io.Reader) (*Mbox, error) {
	s := bufio.NewScanner(r)
	var eof bool
	if s.Scan() { // skip leading From_ line
		// ensure it was a From_ line we skipped
		if !isFromLine(s.Text()) {
			return nil, fmt.Errorf("first line does not begin \"From \"")
		}
	} else {
		eof = true
	}
	return &Mbox{scanner: s, eof: eof}, nil
}

// TODO: add Content-Length parsing so as not to mistake unescaped
// From_ lines in the message body as message separators.

// Message satisfies the folder.Folder interface.
func (m *Mbox) Message() (io.Reader, func() error, error) {
	if m.eof {
		return nil, nil, nil
	}

	lineCh := make(chan []byte)
	go func() {
		for {
			eof := !m.scanner.Scan()
			if eof {
				m.eof = true
				break
			}
			line := m.scanner.Text()
			if isFromLine(line) {
				break
			}
			if isEscapedFromLine(line) {
				line = line[1:] // unescape
			}
			lineCh <- []byte(line)
			lineCh <- []byte("\r\n")
		}
		close(lineCh)
	}()
	closer := func() error {
		for _, ok := <-lineCh; ok; {
		}
		return m.scanner.Err()
	}
	return chanrw.NewReader(lineCh), closer, nil
}

func isFromLine(line string) bool {
	// TODO: use stricter parsing; look for sender and date
	return strings.HasPrefix(line, "From ")
}

func isEscapedFromLine(line string) bool {
	// TODO: in "mboxrd" format there can be multiple leading ">" characters
	return strings.HasPrefix(line, ">From ")
}
