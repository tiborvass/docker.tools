package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

const start = `=== BEGIN goroutine stack dump ===`
const end = `=== END goroutine stack dump ===`

func undump(w io.Writer, r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	s, e := bytes.Index(b, []byte(start)), bytes.Index(b, []byte(end))
	if s < 0 || e < 0 {
		return fmt.Errorf("could not find delimiters")
	}
	b[s+len(start)-1] = '"'
	b[e] = '"'
	str := b[s+len(start)-1 : e+1]
	n := len(str)
	for i := 0; i < n; i++ {
		if str[i] == '\n' {
			copy(str[i:], str[i+1:])
			n--
			i--
			//str[n] = 0
			str = str[:n]
		}
	}
	t, err := strconv.Unquote(string(str))
	if err != nil {
		return fmt.Errorf("unquote: %v", err)
	}
	_, err = fmt.Fprint(w, t)
	return err
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Stderr.Write([]byte{'\n'})
	os.Exit(1)
}

func main() {
	var (
		r   io.ReadCloser
		err error
	)
	switch len(os.Args) {
	case 1:
		r = os.Stdin
	case 2:
		r, err = os.Open(os.Args[1])
		if err != nil {
			fatal(err)
		}
	default:
		fatal(fmt.Errorf("only accept 0 or 1 argument"))
	}
	if err := undump(os.Stdout, r); err != nil {
		fatal(err)
	}
}
