package mailrucloud

import (
	"io"
)

type ioProgress struct {
	r    io.Reader
	max  int
	sent int
	ch   chan<- int
}

func NewIoProgress(r io.Reader, m int, ch chan<- int) *ioProgress {
	if ch != nil {
		ch <- int(m)
	}
	return &ioProgress{r, m, 0, ch}
}

func (i *ioProgress) Read(p []byte) (n int, err error) {
	n, err = i.r.Read(p)
	if i.ch != nil {
		i.sent += n
		i.ch <- i.sent
	}
	return
}
