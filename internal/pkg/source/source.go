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
	Seeker   *os.File
	reader   *bufio.Reader
	tempFile *os.File
	offset   int64
}

func (is *Source) Close() (err error) {
	if is.tempFile != nil {
		err = is.tempFile.Close()
	}
	e := is.Seeker.Close()
	if e != nil {
		err = e
	}

	return err
}

func File(input *os.File, cfg *config.Config) (*Source, error) {
	var err error
	result := &Source{}

	// If it's file we can open it again for seeking.
	result.Seeker, err = os.Open(input.Name())
	if err != nil {
		return nil, err
	}

	reader := io.LimitReader(input, cfg.MaxFileSizeBytes)
	result.reader = bufio.NewReaderSize(reader, maxLineSize)

	return result, nil
}

func Reader(input io.Reader, cfg *config.Config) (*Source, error) {
	var err error
	result := &Source{}

	// We will write the as read to a temp file.  Seek against the temp file.
	result.tempFile, err = os.CreateTemp("", "jvl-*.log")
	if err != nil {
		return nil, err
	}
	reader := io.TeeReader(input, result.tempFile)

	result.Seeker, err = os.Open(result.tempFile.Name())
	if err != nil {
		result.tempFile.Close()

		return nil, err
	}

	reader = io.LimitReader(reader, cfg.MaxFileSizeBytes)
	result.reader = bufio.NewReaderSize(reader, maxLineSize)

	return result, nil
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

// ReadLogEntry reads the next ReadLogEntry from the file.
func (is *Source) ReadLogEntry() (LazyLogEntry, error) {
	for {
		line, err := is.reader.ReadSlice(byte('\n'))
		if err != nil {
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
