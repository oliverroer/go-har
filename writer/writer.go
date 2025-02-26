package harwriter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/oliverroer/go-har"
)

type EntryWriter struct {
	name    string
	file    *os.File
	encoder *json.Encoder
}

func DefaultName() string {
	now := time.Now().UTC()
	const format = "2006-01-02_15-04-05.000000000"
	name := now.Format(format)
	return name
}

func Open(name string) (*EntryWriter, error) {
	cleaned := filepath.Clean(name)
	file, err := os.Create(cleaned)
	if err != nil {
		return nil, err
	}

	writer := EntryWriter{
		name:    name,
		file:    file,
		encoder: json.NewEncoder(file),
	}

	return &writer, nil
}

func (w *EntryWriter) WriteEntry(
	request har.Request,
	response har.Response,
	startedAt time.Time,
	duration time.Duration,
) error {
	entry := har.Entry{
		StartedDateTime: startedAt,
		Time:            int(duration.Milliseconds()),
		Request:         request,
		Response:        response,
	}

	err := w.encoder.Encode(entry)
	if err != nil {
		return err
	}

	_ = w.file.Sync()

	return nil
}

func (w *EntryWriter) Close() error {
	return w.file.Close()
}

func (w *EntryWriter) RoundTripper(base http.RoundTripper) http.RoundTripper {
	return &harRoundTripper{
		base:   base,
		writer: w,
	}
}

const harBegin = `{
  "log": {
    "version": "1.2",
    "creator": {
      "name": "github.com/oliverroer/go-har",
      "version": "0.1.1"
    },
    "entries": [`

const harEnd = `
    ]
  }
}
`

func EntriesToHar(harFile string, entryFiles ...string) error {
	harFile = filepath.Clean(harFile)
	har, err := os.Create(harFile)
	if err != nil {
		return err
	}

	_, err = har.WriteString(harBegin)
	if err != nil {
		return err
	}

	for _, entryFile := range entryFiles {
		entryFile = filepath.Clean(entryFile)
		file, err := os.Open(entryFile)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			str := fmt.Sprintf("\n      %s,", line)
			_, err := har.WriteString(str)
			if err != nil {
				return err
			}
		}
	}

	// step one byte back to discard the last trailing comma
	_, err = har.Seek(-1, 1)
	if err != nil {
		return err
	}

	_, err = har.WriteString(harEnd)
	if err != nil {
		return err
	}

	return har.Close()
}
