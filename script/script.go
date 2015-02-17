package script

import (
	"fmt"
	"io"
	"log"
	"os"
)

func NewFilesRecoding(fileprefix string, stdin io.Reader, stdout io.Writer, stderr io.Writer) (in io.Reader, out io.Writer, err io.Writer, fail error) {
	if fd, fail := os.Create(fileprefix + ".in"); fail == nil {
		in = NewReader(stdin, fd)
	} else {
		log.Panic(fail)
	}

	if fd, fail := os.Create(fileprefix + ".out"); fail == nil {
		out = NewWriter(stdout, fd)
	} else {
		log.Panic(fail)
	}

	if fd, fail := os.Create(fileprefix + ".err"); fail == nil {
		err = NewWriter(stderr, fd)
	} else {
		log.Panic(fail)
	}
	return
}

func NewFilesTimingRecoding(fileprefix string, stdin io.Reader, stdout io.Writer, stderr io.Writer) (in io.Reader, out io.Writer, err io.Writer, fail error) {
	in, out, err, fail = NewFilesRecoding(fileprefix, stdin, stdout, stderr)

	if in, ok := in.(TickInfo); ok {
		if fail := makeTicker(fileprefix+".in", in); fail != nil {
			log.Println(fail)
		}
	} else {
		log.Println("No ok cast...")
	}

	if out, ok := in.(TickInfo); ok {
		if fail := makeTicker(fileprefix+".out", out); fail != nil {
			log.Println(fail)
		}
	} else {
		log.Println("No ok cast...")
	}

	if err, ok := in.(TickInfo); ok {
		if fail := makeTicker(fileprefix+".err", err); fail != nil {
			log.Println(fail)
		}
	} else {
		log.Println("No ok cast...")
	}

	return
}

func makeTicker(fileprefix string, ticker TickInfo) (err error) {
	var fd *os.File
	fd, err = os.Create(fileprefix + ".time")
	if err != nil {
		return
	}

	go func() {
		c := ticker.Ticker()
		log.Print("time record")
		for {
			p, ok := <-c
			if !ok {
				break
			}
			fmt.Fprintln(fd, p)
			fd.Sync()
		}
		fd.Close()
	}()
	return nil
}
