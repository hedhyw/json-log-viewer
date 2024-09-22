package source

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

const (
	// initialLogSize is a capacity of a slice with logs.
	initialLogSize int = 1000

	// RefreshInterval is an interval of refreshing logs.
	RefreshInterval = 200 * time.Millisecond
)

// StartStreaming synchronizes log entries with the file and sends them to the channel.
func (s *Source) StartStreaming(ctx context.Context, send func(msg LazyLogEntries, err error)) {
	logEntriesLock := sync.Mutex{}
	logEntries := make([]LazyLogEntry, 0, initialLogSize)
	eofEvent := make(chan struct{}, 1)

	// Load log entries async..
	go s.readLogEntries(ctx, send, &logEntriesLock, &logEntries, eofEvent)

	// periodically send new log entries to the program.
	go func() {
		ticker := time.NewTicker(RefreshInterval)
		lastLen := -1
		defer ticker.Stop()

		sendUpdates := func() {
			// Only send log update the program state every ticker seconds,
			// to avoid stressing the main loop.
			logEntriesLock.Lock()
			logEntriesClone := make([]LazyLogEntry, len(logEntries))
			copy(logEntriesClone, logEntries)
			logEntriesLock.Unlock()

			nextLen := len(logEntriesClone)
			if lastLen != nextLen {
				send(LazyLogEntries{
					Seeker:  s.Seeker,
					Entries: logEntriesClone,
				}, nil)
				lastLen = nextLen
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-eofEvent:
				sendUpdates()

				return
			case <-ticker.C:
				sendUpdates()
			}
		}
	}()
}

func (s *Source) readLogEntries(
	ctx context.Context,
	send func(msg LazyLogEntries, err error),
	logEntriesLock *sync.Mutex,
	logEntries *[]LazyLogEntry,
	eofEvent chan struct{},
) {
	defer func() {
		eofEvent <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		entry, err := s.readLogEntry()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if !s.CanFollow() {
					return
				}

				// wait for new log entries to be written to the file,
				// and try again.
				ticker := time.NewTicker(RefreshInterval)
				select {
				case <-ctx.Done():
					ticker.Stop()
					return
				case <-ticker.C:
					ticker.Stop()
				}

				continue
			}
			send(LazyLogEntries{}, err)
			return
		}

		logEntriesLock.Lock()
		*logEntries = append(*logEntries, entry)
		logEntriesLock.Unlock()
	}
}
