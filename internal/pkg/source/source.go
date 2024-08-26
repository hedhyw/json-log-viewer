package source

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

const (
	maxLineSize = 8 * 1024 * 1024
)

type Source struct {
	Seeker         *os.File
	reader         *bufio.Reader
	file           *os.File
	offset         int64
	prevFollowSize int64
	name           string
	maxSize        int64
}

func (is *Source) Close() (err error) {
	err = is.file.Close()
	e := is.Seeker.Close()
	if e != nil {
		err = e
	}
	return err
}

func File(name string, cfg *config.Config) (*Source, error) {
	var err error
	is := &Source{
		maxSize: cfg.MaxFileSizeBytes,
		name:    name,
	}

	is.file, err = os.Open(name)
	if err != nil {
		return nil, err
	}

	is.Seeker, err = os.Open(name)
	if err != nil {
		_ = is.file.Close()
		return nil, err
	}

	is.reader = bufio.NewReaderSize(io.LimitReader(is.file, is.maxSize), maxLineSize)
	return is, nil
}

func Reader(input io.Reader, cfg *config.Config) (*Source, error) {
	var err error
	is := &Source{
		maxSize: cfg.MaxFileSizeBytes,
	}

	// We will write the as read to a temp file.  Seek against the temp file.
	is.file, err = os.CreateTemp("", "jvl-*.log")
	if err != nil {
		return nil, err
	}

	reader := io.TeeReader(input, is.file)

	is.Seeker, err = os.Open(is.file.Name())
	if err != nil {
		_ = is.file.Close()
		return nil, err
	}

	reader = io.LimitReader(reader, is.maxSize)
	is.reader = bufio.NewReaderSize(reader, maxLineSize)
	return is, nil
}

func (is *Source) ParseLogEntries() (LazyLogEntries, error) {
	logEntries := make([]LazyLogEntry, 0, 1000)
	for {
		entry, err := is.ReadLogEntry()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return LazyLogEntries{}, err
		}
		logEntries = append(logEntries, entry)
	}

	return LazyLogEntries{
		Seeker:  is.Seeker,
		Entries: logEntries,
	}, nil
}

func (is *Source) CanFollow() bool {
	return len(is.name) != 0
}

var FileTruncated = errors.New("file truncated")

// ReadLogEntry reads the next ReadLogEntry from the file.
func (is *Source) ReadLogEntry() (LazyLogEntry, error) {
	for {
		if is.reader == nil {
			// If we can't follow the file, or we have reached the max size, we are done.
			if !is.CanFollow() || is.offset >= is.maxSize {
				return LazyLogEntry{}, io.EOF
			}

			// has the file size changed since we last looked?
			info, err := os.Stat(is.name)
			if err == nil && is.prevFollowSize != info.Size() {
				if info.Size() < is.offset {
					// the file has been truncated or rolled over, all previous line offsets are invalid.
					// we can't recover from this.
					return LazyLogEntry{}, FileTruncated
				}
				is.prevFollowSize = info.Size()
				// reset the reader and try to read the file again.
				_, _ = is.file.Seek(is.offset, io.SeekStart)
				is.reader = bufio.NewReaderSize(io.LimitReader(is.file, is.maxSize-is.offset), maxLineSize)
			} else {
				return LazyLogEntry{}, io.EOF
			}
		}

		line, err := is.reader.ReadSlice(byte('\n'))
		if err != nil {
			if errors.Is(err, io.EOF) {
				// set the reader to nil so that we can recover from EOF.
				is.reader = nil
			}
			return LazyLogEntry{}, err
		}
		length := len(line)
		offset := is.offset
		is.offset += int64(length)
		if len(bytes.TrimSpace(line)) != 0 {
			return LazyLogEntry{
				offset: offset,
				length: length,
			}, nil
		}
	}
}
