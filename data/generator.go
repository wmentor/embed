// +build ignore

package main

import (
	"github.com/wmentor/embed"
)

func main() {

	err := embed.Make("testdata.txt", "testdata.txt.go", "data", "github.com/wmentor/embed/data/testdata.txt")
	if err != nil {
		panic(err)
	}

}
