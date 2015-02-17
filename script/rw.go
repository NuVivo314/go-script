package script

import (
	"io"
	"log"
	"os"
	"time"
)

type TickInfo interface {
	Ticker() chan *Tick
}

type read struct {
	r       io.Reader
	w       *os.File
	t       time.Time
	ticking bool
	tick    chan *Tick
}

func NewReader(r io.Reader, w *os.File) io.Reader {
	return &read{r, w, time.Now(), false, nil}
}

func (tee *read) Read(p []byte) (n int, err error) {
	n, err = tee.r.Read(p)
	tee.update(n)
	if n > 0 {
		if n, err := tee.w.Write(p[:n]); err != nil {
			return n, err
		}
		tee.w.Sync()
	}
	return
}

func (tee *read) Close() error {
	tee.w.Close()
	if tee.tick != nil {
		close(tee.tick)
	}
	return nil
}

func (tee *read) Ticker() chan *Tick {
	if tee.tick == nil {
		tee.tick = make(chan *Tick, 5)
	}

	tee.ticking = true
	tee.t = time.Now()

	return tee.tick
}

func (tee *read) update(n int) {
	if !tee.ticking {
		return
	}
	newT := time.Now()
	tee.tick <- &Tick{Since(tee.t), n}
	tee.t = newT
}

// &&&
type write struct {
	in      io.Writer
	out     *os.File
	t       time.Time
	ticking bool
	tick    chan *Tick
}

func NewWriter(in io.Writer, out *os.File) io.WriteCloser {
	return &write{in, out, time.Now(), false, nil}
}

func (pipe *write) Write(data []byte) (n int, err error) {
	n, err = pipe.in.Write(data)
	pipe.update(n)
	if pipe.out != nil {
		if _, err = pipe.out.Write(data[:n]); err != nil {
			log.Println(err)
			pipe.out = nil
		}
		pipe.out.Sync()
	}
	return
}

func (pipe *write) Close() error {
	pipe.out.Close()
	if cout, ok := pipe.in.(io.Closer); ok {
		cout.Close()
	}
	if pipe.tick != nil {
		close(pipe.tick)
	}

	return nil
}

func (pipe *write) Ticker() chan *Tick {
	if pipe.tick == nil {
		pipe.tick = make(chan *Tick, 5)
	}

	pipe.ticking = true
	pipe.t = time.Now()

	return pipe.tick
}

func (pipe *write) update(n int) {
	if !pipe.ticking {
		return
	}
	newT := time.Now()
	pipe.tick <- &Tick{Since(pipe.t), n}
	pipe.t = newT
}
