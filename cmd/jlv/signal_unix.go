//go:build !windows

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// interruptProcessGroup sends SIGINT to the process group of the current
// process.
//
// A shell pipeline (`some-command | jlv`) runs all its commands in a single
// process group, so this stops the command that writes to the standard
// input. It is normally done by the terminal itself on Ctrl+C, but the
// application reads the key in the raw mode, so the terminal doesn't send
// anything.
func interruptProcessGroup() error {
	return interruptGroup(os.Getpid(), syscall.Getpgrp(), signal.Ignore, syscall.Kill)
}

func interruptGroup(
	pid int,
	pgid int,
	ignore func(signals ...os.Signal),
	kill func(pid int, signal syscall.Signal) error,
) error {
	// The group contains only the application itself and its children,
	// there is nothing to notify.
	if pid == pgid {
		return nil
	}

	// The application is already exiting, so it should not interrupt itself.
	ignore(os.Interrupt)

	// Zero means the process group of the current process.
	if err := kill(0, syscall.SIGINT); err != nil {
		return fmt.Errorf("signaling the process group: %w", err)
	}

	return nil
}
