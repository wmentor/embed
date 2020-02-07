package embed

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	data map[string][]byte

	ErrNotFound       error
	ErrPackageName    error
	ErrReadSourceFile error
	ErrWriteFile      error
)

func init() {
	data = map[string][]byte{}

	ErrNotFound = errors.New("embed resource not found")
	ErrPackageName = errors.New("invalid package name")
	ErrReadSourceFile = errors.New("read source file error")
	ErrWriteFile = errors.New("error write file")
}

func Register(name, content string) {

	cont, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		panic(err)
	}

	gz, _ := gzip.NewReader(bytes.NewReader(cont))

	c1, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	data[name] = c1
}

func Delete(name string) {
	delete(data, name)
}

func Get(name string) (io.Reader, error) {

	if val, has := data[name]; has {
		return bytes.NewReader(val), nil
	}

	return nil, ErrNotFound
}

func Make(src string, dest string, pkg string, name string) error {
	if pkg == "" {
		return ErrPackageName
	}

	if src == "" {
		return ErrReadSourceFile
	}

	if dest == "" {
		return ErrWriteFile
	}

	data, err := ioutil.ReadFile(src)
	if err != nil {
		return ErrReadSourceFile
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	if _, err := gz.Write(data); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	data = b.Bytes()

	rw, err := os.Create(dest)
	if err != nil {
		return ErrWriteFile
	}
	defer rw.Close()

	result := base64.StdEncoding.EncodeToString(data)

	builder := strings.Builder{}

	fmt.Fprintf(rw, "package %s\n\nimport (\n\t\"github.com/wmentor/embed\"\n)\n\nfunc init() {\n\n\tembed.Register(%s, `\n", pkg, strconv.Quote(name))

	i := 0

	for _, r := range result {
		i++
		builder.WriteRune(r)
		if i >= 64 {
			builder.WriteRune('\n')
			fmt.Fprint(rw, builder.String())
			i = 0
			builder.Reset()
		}
	}

	if i > 0 {
		fmt.Fprint(rw, builder.String())
	}

	fmt.Fprintf(rw, "`)\n\n}\n")

	return nil
}
