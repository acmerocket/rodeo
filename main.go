package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func parse_args() []string {
	// wordPtr := flag.String("word", "foo", "a string")
	// numbPtr := flag.Int("numb", 42, "an int")
	// forkPtr := flag.Bool("fork", false, "a bool")
	// var svar string
	// flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()
	return flag.Args()
}

func cleanup() {
	type_report()
}

func main() {
	// handle the sigterm
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	// parse cli args
	names := parse_args()

	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)

	// create a new renderer
	renderer, err := build_renderer(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// fill the buffer
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err) // FIXME
		}

		// output the buffer
		render_buffer(buf, names, renderer)
		time.Sleep(1 * time.Millisecond) // self throttle
	}
}
