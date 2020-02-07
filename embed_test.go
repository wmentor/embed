package embed

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestEmbed(t *testing.T) {

	wait := []byte(`Hello, World!
It's my test data.`)

	Register("github.com/wmentor/embed/data/testdata.txt", `
H4sIAAAAAAAA//JIzcnJ11EIzy/KSVHk8ixRL1bIrVQoSS0uUUhJLEnUAwQAAP//
YIsh6yAAAAA=`)

	in, err := Get("github.com/wmentor/embed/data/testdata.txt")
	if err != nil {
		t.Fatal("Data not found")
	}

	data, _ := ioutil.ReadAll(in)

	if bytes.Compare(data, wait) != 0 {
		t.Fatal("Invalid data")
	}

	in, err = Get("github.com/wmentor")
	if err != ErrNotFound || in != nil {
		t.Fatal("Return data")
	}
}
