package harwriter

import (
	"net/http"
	"time"

	"github.com/oliverroer/go-har"
)

var _ http.RoundTripper = (*harRoundTripper)(nil)

type harRoundTripper struct {
	base   http.RoundTripper
	writer *HarWriter
}

func (t *harRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	harRequest := har.RequestFromHttpRequest(req)

	start := time.Now()
	res, err := t.base.RoundTrip(req)
	elapsed := time.Since(start)

	harResponse := har.ResponseFromHttpResponse(res)

	_ = t.writer.WriteEntry(harRequest, harResponse, start, elapsed)

	return res, err
}
