package har

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

func RequestFromHttpRequest(req *http.Request) Request {
	headerData := headerDataFromHttpHeader(req.Header)

	request := Request{
		Method:      req.Method,
		URL:         req.URL.String(),
		HTTPVersion: req.Proto,
		Cookies:     []Cookie{}, // TODO: Implement cookies?
		Headers:     headerData.headers,
		QueryString: []QueryString{}, // TODO: Implement query string?
		PostData:    nil,
		HeadersSize: headerData.size,
		BodySize:    -1,
	}

	if body := peekBody(&req.Body); body != nil {
		request.PostData = &PostData{
			MimeType: headerData.mimeType,
			Text:     string(body),
		}
		request.BodySize = len(body)
	}

	return request
}

func ResponseFromHttpResponse(res *http.Response) Response {
	headerData := headerDataFromHttpHeader(res.Header)

	response := Response{
		Status:      res.StatusCode,
		StatusText:  res.Status,
		HttpVersion: res.Proto,
		Cookies:     []Cookie{}, // TODO: Implement cookies?
		Headers:     headerData.headers,
		Content: Content{
			MimeType: headerData.mimeType,
			Size:     -1,
		},
		RedirectURL: headerData.location,
		HeadersSize: headerData.size,
		BodySize:    int(res.ContentLength),
	}

	if body := peekBody(&res.Body); body != nil {
		response.Content.Text = string(body)
		response.Content.Size = len(body)
	}

	return response
}

type headerData struct {
	headers  []Header
	size     int
	mimeType string
	location string
}

func headerDataFromHttpHeader(httpHeader http.Header) headerData {
	headerCount := len(httpHeader)
	data := headerData{
		headers:  make([]Header, 0, headerCount),
		size:     0,
		mimeType: "",
		location: "",
	}

	if headerCount == 0 {
		return data
	}

	for name, values := range httpHeader {
		value := strings.Join(values, " ")

		header := Header{
			Name:  name,
			Value: value,
		}

		data.headers = append(data.headers, header)

		// Count
		// - the length of the header name
		// - the length of ": "
		// - the length of the header value
		// - the length of CRLF
		data.size += len(name) + 2 + len(value) + 2

		switch name {
		case "Content-Type":
			data.mimeType = value

		case "Location":
			data.location = value
		}
	}

	// Count the additional CRLF
	data.size += 2

	return data
}

// peekBody reads all bytes from an io.ReadCloser without "consuming" it,
// by replacing it with a new io.ReadCloser that can be read again.
//
// This is heavily inspired by drainBody which is used internally in various
// dump functions in the httputil package.
func peekBody(body *io.ReadCloser) []byte {
	readCloser := *body
	if readCloser == nil || readCloser == http.NoBody {
		return nil
	}

	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(readCloser); err != nil {
		return nil
	}
	if err := readCloser.Close(); err != nil {
		return nil
	}

	peeked := buffer.Bytes()

	reader := bytes.NewReader(peeked)
	*body = io.NopCloser(reader)

	return peeked
}
