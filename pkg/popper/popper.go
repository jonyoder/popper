package popper

import "io"

// Example
//
// 1. Obtain a reader that needs to write somewhere, e.g. `resp.Body`
// 2. Create a popper with:
//      pr := NewPopper(resp.Body)
// 3. Write to a file with:
//      io.Copy(dst, pr)
// 4. A new consumer needs to read while still writing:
//      r := pr.NewReader()
//      io.Copy(dst, r)
//

type Popper interface {
	io.Reader
	NewReader(r io.Reader) io.Reader
}

type PoppingReader interface {
	io.Reader
}

type DefaultPopper struct {
	io.Reader
}

func NewPopper(r io.Reader) *DefaultPopper {
	return &DefaultPopper{
		Reader: r,
	}
}

func (p *DefaultPopper) NewReader(r io.Reader) io.Reader {
	return &poppingReader{
		Reader: r,
	}
}

type poppingReader struct {
	io.Reader
}

func (pr *poppingReader) Read(p []byte) (n int, err error) {
	return pr.Reader.Read(p)
}
