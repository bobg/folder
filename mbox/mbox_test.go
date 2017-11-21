package mbox

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	msg1 = `Subject: message 1

body
>From foo bar
still the body
`

	msg2 = `Subject: message 2

body
`
)

var f = fmt.Sprintf("From foo bar\n%s\nFrom foo bar\n%s\n", msg1, msg2)

func TestMbox(t *testing.T) {
	m, err := New(strings.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	msg1R, closer, err := m.Message()
	if err != nil {
		t.Fatalf("getting message 1: %s", err)
	}
	msg1Bytes, err := ioutil.ReadAll(msg1R)
	if err != nil {
		t.Fatalf("reading message 1: %s", err)
	}
	if !sameMsg(string(msg1Bytes), msg1) {
		t.Errorf("message 1: got:\n%s\nwant:\n%s", string(msg1Bytes), msg1)
	}
	err = closer()
	if err != nil {
		t.Fatalf("closing message 1: %s", err)
	}
	msg2R, closer, err := m.Message()
	if err != nil {
		t.Fatalf("getting message 2: %s", err)
	}
	msg2Bytes, err := ioutil.ReadAll(msg2R)
	if err != nil {
		t.Fatalf("reading message 2: %s", err)
	}
	if !sameMsg(string(msg2Bytes), msg2) {
		t.Errorf("message 2: got:\n%s\nwant:\n%s", string(msg2Bytes), msg2)
	}
	err = closer()
	if err != nil {
		t.Fatalf("closing message 2: %s", err)
	}
	msg3R, closer, err := m.Message()
	if err != nil {
		t.Fatalf("getting message 3: %s", err)
	}
	if msg3R != nil {
		t.Error("got non-nil reader for message 3")
	}
}

func sameMsg(m1, m2 string) bool {
	s1 := bufio.NewScanner(strings.NewReader(m1))
	s2 := bufio.NewScanner(strings.NewReader(m2))

	for {
		ok1 := s1.Scan()
		ok2 := s2.Scan()

		if ok1 && ok2 {
			l1 := s1.Text()
			l2 := s2.Text()

			if isEscapedFromLine(l2) {
				l2 = l2[1:]
			}

			if l1 != l2 {
				return false
			}
			continue
		}

		if !ok1 && !ok2 {
			return true
		}

		if ok1 {
			return onlyEmptyLines(s1)
		}
		return onlyEmptyLines(s2)
	}
}

func onlyEmptyLines(s *bufio.Scanner) bool {
	for {
		// s.Scan() has already been called
		if s.Text() != "" {
			return false
		}
		if !s.Scan() {
			return true
		}
	}
}
