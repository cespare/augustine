package augustine

import (
	"bufio"
	"io"
)

type TokenReader interface {
	ReadToken() (tok []byte, err error)
}

type Tokenizer func(io.Reader) TokenReader

type SimpleTokenReader struct {
	s *bufio.Scanner
}

func SimpleTokenizer(r io.Reader) TokenReader {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)
	return &SimpleTokenReader{s}
}

func (r *SimpleTokenReader) ReadToken() (tok []byte, err error) {
	if r.s.Scan() {
		b := r.s.Bytes()
		tok = make([]byte, len(b))
		copy(tok, b)
		return tok, nil
	}
	if err := r.s.Err(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}
