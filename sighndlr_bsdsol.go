// Copyright © 2020 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// +build freebsd openbsd solaris darwin

// Signal handling for FreeBSD, OpenBSD, Darwin, and Solaris
// systems that have SIGINFO

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/alecthomas/units"
	"golang.org/x/sys/unix"
)

//signalNotifySetup sets up the signals and their channel
func (ac *AgentController) signalNotifySetup() {
	signal.Notify(ac.sig, os.Interrupt, unix.SIGTERM, unix.SIGHUP, unix.SIGPIPE, unix.SIGINFO)
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
			case os.Interrupt, unix.SIGTERM:
				ac.cncl()
				return
				// log.Println("waiting for final metric flushes... press CTRL+C to exit without flushing")
				// break
			case unix.SIGPIPE, unix.SIGHUP:
				// Noop
			case unix.SIGINFO:
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
