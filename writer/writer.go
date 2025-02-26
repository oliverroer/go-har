package harwriter

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/oliverroer/go-har"
)

type HarWriter struct {
	file       *os.File
	encoder    *json.Encoder
	entryCount int
}

const beginJson = `{
  "log": {
    "version": "1.2",
    "creator": {
      "name": "github.com/oliverroer/go-har",
      "version": "0.1.0"
    },
    "entries": [`

const endJson = `    ]
  }
}
`

const INDENT_1 = "  "
const INDENT_2 = INDENT_1 + INDENT_1
const INDENT_3 = INDENT_1 + INDENT_2
const NEWLINE = "\n"

func DefaultName() string {
	now := time.Now().UTC()
	const format = "2006-01-02_15-04-05.000000000"
	name := now.Format(format)
	return name + ".har"
}

func OpenDefault(dir string) (*HarWriter, error) {
	perm := fs.FileMode(0750)
	err := os.MkdirAll(dir, perm)
	if err != nil {
		return nil, err
	}

	name := path.Join(dir, DefaultName())
	return Open(name)
}

func Open(name string) (*HarWriter, error) {
	cleaned := filepath.Clean(name)
	file, err := os.Create(cleaned)
	if err != nil {
		return nil, err
	}

	writer := HarWriter{
		file:       file,
		encoder:    json.NewEncoder(file),
		entryCount: 0,
	}
	writer.encoder.SetIndent(INDENT_3, INDENT_1)

	_, err = file.WriteString(beginJson)
	if err != nil {
		return nil, err
	}

	return &writer, nil
}

func (w *HarWriter) WriteEntry(
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

	if w.entryCount > 0 {
		// previous encode call inserted a premature newline
		// so we seek back and replace it with a comma
		_, err := w.file.Seek(-1, 1)
		if err != nil {
			return err
		}
		_, err = w.file.WriteString(",")
		if err != nil {
			return err
		}
	}
	_, err := w.file.WriteString(NEWLINE + INDENT_3)
	if err != nil {
		return err
	}

	err = w.encoder.Encode(entry)
	if err != nil {
		return err
	}
	w.entryCount += 1

	_ = w.file.Sync()

	return nil
}

func (w *HarWriter) Close() error {
	_, err := w.file.WriteString(endJson)
	if err != nil {
		return err
	}

	return w.file.Close()
}

func (w *HarWriter) RoundTripper(base http.RoundTripper) http.RoundTripper {
	return &harRoundTripper{
		base:   base,
		writer: w,
	}
}
