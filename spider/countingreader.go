// Author Daniel Schlaug
// Written at Hong Kong University of Science and Technology in 2013

package spider

import (
	"io"
)

type CountingReader struct {
	reader io.Reader
	count  int64
}

func NewCountingReader(reader io.Reader) *CountingReader {
	return &CountingReader{
		reader,
		0,
	}
}

func (c *CountingReader) Read(p []byte) (n int, err error) {
	n, err = c.reader.Read(p)
	c.count += int64(n)
	return n, err
}

func (c *CountingReader) ByteCount() int64 {
	return c.count
}
