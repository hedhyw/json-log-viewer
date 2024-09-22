package source

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

const (
	maxLineSize = 8 * 1024 * 1024

	temporaryFilePattern = "jvl-*.log"
)

type Source struct {
	// Seeker is used to do random access reads from the file.
	Seeker *os.File
	// Reader is used to read the file sequentially.
	reader *bufio.Reader
	// The log file we are reading from, or a temp file we are writing to (depending on if created with File or Reader func).
	file *os.File
	// offset is the next offset a long entry will be read from.
	offset int64
	// prevFollowSize is the size of the file the last time we checked
	prevFollowSize int64
	// name is the name of the file we are reading.
	name string
	// maxSize is the maximum size of the file we will read.
	maxSize int64
}

func (s *Source) Close() error {
	return errors.Join(s.file.Close(), s.Seeker.Close())
}

// File creates a new Source for reading log entries from a file.
func File(name string, cfg *config.Config) (*Source, error) {
	var err error

	source := &Source{
		maxSize: cfg.MaxFileSizeBytes,
		name:    name,
	}

	source.file, err = os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}

	source.Seeker, err = os.Open(name)
	if err != nil {
		return nil, errors.Join(err, source.file.Close())
	}

	source.reader = bufio.NewReaderSize(
		io.LimitReader(source.file, source.maxSize),
		maxLineSize,
	)

	return source, nil
}

// Reader creates a new Source for reading log entries from an io.Reader.
// This will write the input to a temp file, which will be used to seek against.
func Reader(input io.Reader, cfg *config.Config) (*Source, error) {
	var err error

	source := &Source{
		maxSize: cfg.MaxFileSizeBytes,
	}

	// We will write the as read to a temp file.  Seek against the temp file.
	source.file, err = os.CreateTemp(
		"", // Default directory for temporary files.
		temporaryFilePattern,
	)
	if err != nil {
		return nil, fmt.Errorf("creating temporary file: %w", err)
	}

	// The io.TeeReader will write the input to the is.file as it is read.
	reader := io.TeeReader(input, source.file)

	// We can now seek against the data that is read in the input io.Reader.
	source.Seeker, err = os.Open(source.file.Name())
	if err != nil {
		return nil, errors.Join(err, source.file.Close())
	}

	reader = io.LimitReader(reader, source.maxSize)
	source.reader = bufio.NewReaderSize(reader, maxLineSize)

	return source, nil
}

func (s *Source) ParseLogEntries() (LazyLogEntries, error) {
	logEntries := make([]LazyLogEntry, 0, initialLogSize)
	for {
		entry, err := s.readLogEntry()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return LazyLogEntries{}, err
		}

		logEntries = append(logEntries, entry)
	}

	return LazyLogEntries{
		Seeker:  s.Seeker,
		Entries: logEntries,
	}, nil
}

func (s *Source) CanFollow() bool {
	return len(s.name) != 0
}

const ErrFileTruncated semerr.Error = "file truncated"

// ReadLogEntry reads the next ReadLogEntry from the file.
func (s *Source) readLogEntry() (LazyLogEntry, error) {
	for {
		if s.reader == nil {
			// If we can't follow the file, or we have reached the max size, we are done.
			if !s.CanFollow() || s.offset >= s.maxSize {
				return LazyLogEntry{}, io.EOF
			}

			// Has the file size changed since we last looked?
			info, err := os.Stat(s.name)
			if err != nil || s.prevFollowSize == info.Size() {
				return LazyLogEntry{}, io.EOF
			}

			if info.Size() < s.offset {
				// The file has been truncated or rolled over, all previous line
				// offsets are invalid. We can't recover from this.
				return LazyLogEntry{}, ErrFileTruncated
			}

			s.prevFollowSize = info.Size()
			// Reset the reader and try to read the file again.
			_, _ = s.file.Seek(s.offset, io.SeekStart)
			s.reader = bufio.NewReaderSize(io.LimitReader(s.file, s.maxSize-s.offset), maxLineSize)
		}

		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				// Set the reader to nil so that we can recover from EOF.
				s.reader = nil
			}

			return LazyLogEntry{}, err
		}

		length := len(line)
		offset := s.offset
		s.offset += int64(length)

		if len(bytes.TrimSpace(line)) != 0 {
			return LazyLogEntry{
				offset: offset,
				length: length,
			}, nil
		}
	}
}
