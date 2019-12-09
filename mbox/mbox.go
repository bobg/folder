// Package mbox implements parsing of mbox-style mail folders.

package mbox

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// Mbox is a parser for a Unix mbox-style mail folder.
type Mbox struct {
	scanner *bufio.Scanner
	eof     bool
	r       io.Reader
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
	return &Mbox{scanner: s, eof: eof, r: r}, nil
}

// TODO: add Content-Length parsing so as not to mistake unescaped
// From_ lines in the message body as message separators.

// Message satisfies the folder.Folder interface.
func (m *Mbox) Message() (io.ReadCloser, error) {
	if m.eof {
		return nil, nil
	}

	pr, pw := io.Pipe()

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
			io.WriteString(pw, line)
			io.WriteString(pw, "\r\n")
		}
		pw.Close()
	}()
	return &readCloser{
		r: pr,
		m: m,
	}, nil
}

func (m *Mbox) Close() error {
	if c, ok := m.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

type readCloser struct {
	r io.Reader
	m *Mbox
}

func (r *readCloser) Read(b []byte) (int, error) {
	return r.r.Read(b)
}

func (r *readCloser) Close() error {
	io.Copy(ioutil.Discard, r.r) // flush
	return r.m.scanner.Err()
}

func isFromLine(line string) bool {
	// TODO: use stricter parsing; look for sender and date
	return strings.HasPrefix(line, "From ")
}

func isEscapedFromLine(line string) bool {
	// TODO: in "mboxrd" format there can be multiple leading ">" characters
	return strings.HasPrefix(line, ">From ")
}
