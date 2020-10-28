// Copyright © 2020 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Signal handling for Windows
// doesn't have SIGINFO, attempt to use SIGTRAP instead...

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/alecthomas/units"
)

//SignalNotifySetup sets up the signals and their channel
func SignalNotifySetup(ch chan os.Signal) {
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGPIPE, syscall.SIGTRAP)
}

//HandleSignals handles exiting the program based on different signals
func HandleSignals(cancel context.CancelFunc, ch chan os.Signal) {
	const stackTraceBufferSize = 1 * units.MiB

	//pre-allocate a buffer for stacktrace
	buf := make([]byte, stackTraceBufferSize)

	for {
		select {
		case sig := <-ch:
			log.Printf("signal %s received\n", sig.String())
			switch sig {
			case os.Interrupt, syscall.SIGTERM:
				cancel()
				log.Println("waiting for final metric flushes... press CTRL+C to exit without flushing")
				break
			case syscall.SIGPIPE, syscall.SIGHUP:
				// Noop
			case syscall.SIGTRAP:
				stackLen := runtime.Stack(buf, true)
				fmt.Printf("=== received SIGINFO ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stackLen])
			default:
				log.Printf("signal %s unsupported", sig.String())
			}
		}
	}
}
