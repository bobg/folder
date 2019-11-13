package folder

import "testing"

// Regression test for an infinite-loop bug.
func TestEmptyAlternative(t *testing.T) {
	f, err := Open("testdata/empty-alternative.v7")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	msgR, err := f.Message()
	if err != nil {
		t.Fatal(err)
	}
	defer msgR.Close()
}
