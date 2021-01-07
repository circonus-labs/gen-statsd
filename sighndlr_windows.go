// Copyright © 2020 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Signal handling for Windows
// doesn't have SIGINFO, attempt to use SIGTRAP instead...

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/alecthomas/units"
)

//signalNotifySetup sets up the signals and their channel
func (ac *AgentController) signalNotifySetup() {
	signal.Notify(ac.sig, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGPIPE, syscall.SIGTRAP)
}

//handleSignals handles exiting the program based on different signals
func (ac *AgentController) handleSignals() {
	const stackTraceBufferSize = 1 * units.MiB

	//pre-allocate a buffer for stacktrace
	buf := make([]byte, stackTraceBufferSize)

	for {
		select {
		case sig := <-ac.sig:
			log.Printf("signal %s received\n", sig.String())
			switch sig {
			case os.Interrupt, syscall.SIGTERM:
				ac.cncl()
				return
			case syscall.SIGPIPE, syscall.SIGHUP:
				// Noop
			case syscall.SIGTRAP:
				stackLen := runtime.Stack(buf, true)
				fmt.Printf("=== received SIGINFO ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stackLen])
			default:
				log.Printf("signal %s unsupported", sig.String())
			}
		case <-ac.ctx.Done():
			return
		}
	}
}
