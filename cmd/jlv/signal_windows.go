//go:build windows

package main

// interruptProcessGroup is not supported on windows, there are no process
// groups that are shared by the commands of a pipeline.
func interruptProcessGroup() error {
	return nil
}
