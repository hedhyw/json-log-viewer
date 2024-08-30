package source

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

const InitialLogSize int = 1000

func (is *Source) StartStreaming(ctx context.Context, send func(msg LazyLogEntries, err error)) {
	logEntriesLock := sync.Mutex{}
	logEntries := make([]LazyLogEntry, 0, InitialLogSize)
	eofEvent := make(chan struct{}, 1)

	// Load log entries async..
	go is.readLogEntries(ctx, send, &logEntriesLock, logEntries, eofEvent)

	// periodically send new log entries to the program.
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		lastLen := -1
		defer ticker.Stop()

		sendUpdates := func() {
			// Only send log update the program state every ticker seconds,
			// to avoid stressing the main loop.
			logEntriesLock.Lock()
			entries := logEntries
			logEntriesLock.Unlock()

			nextLen := len(entries)
			if lastLen != nextLen {
				send(LazyLogEntries{
					Seeker:  is.Seeker,
					Entries: entries,
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

func (is *Source) readLogEntries(ctx context.Context, send func(msg LazyLogEntries, err error), logEntriesLock *sync.Mutex, logEntries []LazyLogEntry, eofEvent chan struct{}) {
	defer func() {
		eofEvent <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		entry, err := is.ReadLogEntry()
		if err != nil {
			if errors.Is(err, io.EOF) {

				if !is.CanFollow() {
					return
				}

				// wait for new log entries to be written to the file,
				// and try again.
				ticker := time.NewTicker(200 * time.Millisecond)
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
		logEntries = append(logEntries, entry)
		logEntriesLock.Unlock()
	}
}
