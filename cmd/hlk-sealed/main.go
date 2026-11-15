package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/number571/hidden-lake/internal/kernel/pkg/sealed"
)

const (
	cExitTimeout = time.Second

	cExitOK    = 0
	cExitError = 1
	cExitBoot  = 2
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// stop must not take longer than cExitTimeout after EOF of stdin or a signal
	forceExit := &sync.Once{}
	startExitTimer := func() {
		forceExit.Do(func() { time.AfterFunc(cExitTimeout, func() { os.Exit(cExitOK) }) })
	}
	go func() {
		<-ctx.Done()
		startExitTimer()
	}()

	err := sealed.Run(ctx, &sEOFReader{fReader: os.Stdin, fOnEOF: startExitTimer}, os.Stdout, os.Stderr)
	if err == nil {
		return cExitOK
	}

	// errors of the package do not contain secrets of the boot line
	_, _ = fmt.Fprintf(os.Stderr, "hlk-sealed: %v\n", err)
	if errors.Is(err, sealed.ErrBoot) {
		return cExitBoot
	}
	return cExitError
}

var (
	_ io.Reader = &sEOFReader{}
)

type sEOFReader struct {
	fReader io.Reader
	fOnEOF  func()
}

func (p *sEOFReader) Read(pB []byte) (int, error) {
	n, err := p.fReader.Read(pB)
	if err != nil {
		p.fOnEOF()
	}
	return n, err
}
