package script

import (
	"fmt"
	"time"
)

type duration time.Duration

func Since(t time.Time) duration {
	return duration(time.Now().Sub(t))
}

func (d duration) String() string {
	return fmt.Sprintf("%.6f", float32(d)/float32(time.Second))
}

type Tick struct {
	d duration
	n int
}

func (t *Tick) String() string {
	return fmt.Sprintf("%s %d", t.d, t.n)
}
