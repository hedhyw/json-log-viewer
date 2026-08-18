//go:build !windows

package main

import (
	"os"
	"syscall"
	"testing"

	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterruptGroup(t *testing.T) {
	t.Parallel()

	const (
		pid  = 100
		pgid = 10
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var (
			ignored  []os.Signal
			signaled []int
		)

		err := interruptGroup(pid, pgid,
			func(signals ...os.Signal) { ignored = append(ignored, signals...) },
			func(pid int, signal syscall.Signal) error {
				assert.Equal(t, syscall.SIGINT, signal)
				signaled = append(signaled, pid)

				return nil
			},
		)
		require.NoError(t, err)

		assert.Equal(t, []os.Signal{os.Interrupt}, ignored)
		assert.Equal(t, []int{0}, signaled)
	})

	t.Run("group_leader", func(t *testing.T) {
		t.Parallel()

		err := interruptGroup(pid, pid,
			func(...os.Signal) { t.Fatal("Should not ignore") },
			func(int, syscall.Signal) error {
				t.Fatal("Should not signal")

				return nil
			},
		)
		require.NoError(t, err)
	})

	t.Run("kill_error", func(t *testing.T) {
		t.Parallel()

		err := interruptGroup(pid, pgid,
			func(...os.Signal) {},
			func(int, syscall.Signal) error { return tests.ErrTest },
		)
		require.ErrorIs(t, err, tests.ErrTest)
	})
}
