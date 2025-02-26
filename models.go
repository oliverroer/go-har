package har

import "time"

type HttpArchive struct {
	Log ArchiveLog `json:"log"`
}

/*
 * The structs in this file are based on the objects that are outlined in
 * the following article on HAR:
 * http://www.softwareishard.com/blog/har-12-spec/
 *
 * NOTE:
 * While the struct model the entire spec, there are various pieces of data
 * that we currently do not fill out, due to limitations on what we're able to
 * observe using standard HTTP tooling in Go - for example, we don't provide
 * a granular breakdown of request timings or resolved server IP addresses.
 *
 * This means that the produced HAR file is not as detailed as one might be used
 * to, e.g. when inspecting network requests in Google Chrome.
 */

// ArchiveLog represents the root of exported data.
type ArchiveLog struct {
	// Version is the version number of the format.
	// If empty, string "1.1" is assumed by default.
	Version string `json:"version"`

	// Name and version info of the log creator application.
	Creator Creator `json:"creator"`

	// Name and version info of used browser.
	Browser *Browser `json:"browser,omitempty"`

	// Page is a list of all exported (tracked) pages.
	// Leave out this field if the application does not support grouping by pages.
	Page []Page `json:"page,omitempty"`

	// Entries is a list of all exported (tracked) requests.
	// Sorting entries by startedDateTime (starting from the oldest) is the
	// preferred way to export data since it can make importing faster.
	// However, the reader application should always make sure the array is
	// sorted (if required for the import).
	Entries []Entry `json:"entries"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// Name and version info of the log creator application.
type Creator struct {
	// Name is the name of the application used to export the log.
	Name string `json:"name"`

	// Version is the version of the application used to export the log.
	Version string `json:"version"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// Name and version info of used browser.
type Browser struct {
	// Name is the name of the browser used to export the log.
	Name string `json:"name"`

	// Version is the version of the browser used to export the log.
	Version string `json:"version"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// This object represents an exported page.
type Page struct {
	// StartedDateTime is a date and time stamp for the beginning of the page load
	// (ISO 8601 - YYYY-MM-DDThh:mm:ss.sTZD,
	// e.g. 2009-07-24T19:20:30.45+01:00).
	StartedDateTime time.Time `json:"startedDateTime"`

	// ID is a unique identifier of a page within the log.
	// Entries use it to refer the parent page.
	ID string `json:"id"`

	// Title is the page title.
	Title string `json:"title"`

	// PageTimings has detailed timing info about page load.
	PageTimings PageTimings `json:"pageTimings"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

type PageTimings struct {
	// OnContentLoad is the number of milliseconds since page load started
	// (page.startedDateTime).
	// Use -1 if the timing does not apply to the current request.
	OnContentLoad int `json:"onContentLoad,omitempty"`

	// OnLoad is when the page is loaded (onLoad event fired).
	// Number of milliseconds since page load started (page.startedDateTime).
	// Use -1 if the timing does not apply to the current request.
	OnLoad int `json:"onLoad,omitempty"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// This object represents an exported HTTP request.
type Entry struct {
	// Pageref is a reference to the parent page.
	// Leave out this field if the application does not support grouping by pages.
	Pageref string `json:"pageref,omitempty"`

	// StartedDateTime is the date and time stamp of the request start
	// (ISO 8601 - YYYY-MM-DDThh:mm:ss.sTZD,
	// e.g. 2009-07-24T19:20:30.123+02:00).
	StartedDateTime time.Time `json:"startedDateTime"`

	// Time is the total elapsed time of the request in milliseconds.
	// This is the sum of all timings available in the timings object
	// (i.e. not including -1 values).
	Time int `json:"time"`

	// Request has detailed info about the request.
	Request Request `json:"request"`

	// Response has detailed info about the response.
	Response Response `json:"response"`

	// Cache has info about cache usage.
	Cache Cache `json:"cache"`

	// Timings has detailed timing info about request/response round trip.
	Timings Timings `json:"timings"`

	// ServerIPAddress is the IP address of the server that was connected to
	// (result of DNS resolution).
	ServerIPAddress string `json:"serverIPAddress,omitempty"`

	// Connection is a Unique ID of the parent TCP/IP connection,
	// can be the client or server port number.
	// Note that a port number doesn't have to be unique identifier in cases
	// where the port is shared for more connections.
	// If the port isn't available for the application, any other unique
	// connection ID can be used instead (e.g. connection index).
	// Leave out this field if the application doesn't support this info.
	Connection string `json:"connection,omitempty"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// This object contains detailed info about the performed request.
type Request struct {
	// Request method (GET, POST, ...).
	Method string `json:"method"`

	// Absolute URL of the request (fragments are not included).
	URL string `json:"url"`

	// Request HTTP Version.
	HTTPVersion string `json:"httpVersion"`

	// List of cookie objects.
	Cookies []Cookie `json:"cookies"`

	// List of header objects.
	Headers []Header `json:"headers"`

	// List of query parameter objects.
	QueryString []QueryString `json:"queryString"`

	// Posted data info.
	PostData *PostData `json:"postData,omitempty"`

	// Total number of bytes from the start of the HTTP request message until
	// (and including) the double CRLF before the body.
	// Set to -1 if the info is not available.
	HeadersSize int `json:"headersSize"`

	// Size of the request body (POST data payload) in bytes.
	// Set to -1 if the info is not available.
	BodySize int `json:"bodySize"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// This object contains detailed info about the response.
type Response struct {
	// Response status.
	Status int `json:"status"`

	// Response status description.
	StatusText string `json:"statusText"`

	// Response HTTP Version.
	HttpVersion string `json:"httpVersion"`

	// List of cookie objects.
	Cookies []Cookie `json:"cookies"`

	// List of header objects.
	Headers []Header `json:"headers"`

	// Details about the response body.
	Content Content `json:"content"`

	// Redirection target URL from the Location response header.
	RedirectURL string `json:"redirectURL"`

	// Total number of bytes from the start of the HTTP response message until
	// (and including) the double CRLF before the body.
	// Set to -1 if the info is not available.
	// Note:
	// The size of received response-headers is computed only from headers that
	// are really received from the server.
	// Additional headers appended by the browser are not included in this
	// number, but they appear in the list of header objects.
	HeadersSize int `json:"headersSize"`

	// Size of the received response body in bytes.
	// Set to zero in case of responses coming from the cache (304).
	// Set to -1 if the info is not available.
	BodySize int `json:"bodySize"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// This object contains list of all cookies
// (used in Request and Response objects).
type Cookie struct {
	// The name of the cookie.
	Name string `json:"name"`

	// The cookie value.
	Value string `json:"value"`

	// The path pertaining to the cookie.
	Path string `json:"path"`

	// The host of the cookie.
	Domain string `json:"domain,omitempty"`

	// Cookie expiration time.
	// (ISO 8601 - YYYY-MM-DDThh:mm:ss.sTZD,
	// e.g. 2009-07-24T19:20:30.123+02:00).
	Expires string `json:"expires,omitempty"`

	// Set to true if the cookie is HTTP only, false otherwise.
	HTTPOnly bool `json:"httpOnly,omitempty"`

	// True if the cookie was transmitted over ssl, false otherwise.
	Secure bool `json:"secure,omitempty"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// This object contains details of a header
// (used in Request and Response objects).
type Header struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Comment string `json:"comment,omitempty"`
}

// This object contains describes a value parsed from a query string,
// (embedded in Request object).
type QueryString struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Comment string `json:"comment,omitempty"`
}

// This object describes posted data, if any
// (embedded in Request object).
// Note:
// Text and params fields are mutually exclusive.
type PostData struct {
	// Mime type of posted data.
	MimeType string `json:"mimeType"`

	// List of posted parameters (in case of URL encoded parameters).
	Params []Param `json:"params,omitempty"`

	// Plain text posted data.
	Text string `json:"text"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// List of posted parameters, if any
// (embedded in PostData object).
type Param struct {
	// Name of a posted parameter.
	Name string `json:"name"`

	// Value of a posted parameter or content of a posted file.
	Value string `json:"value,omitempty"`

	// Name of a posted file.
	FileName string `json:"fileName,omitempty"`

	// Content type of a posted file.
	ContentType string `json:"contentType,omitempty"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// This object describes details about response content
// (embedded in Response object).
type Content struct {
	// Length of the returned content in bytes.
	// Should be equal to response.bodySize if there is no compression and
	// bigger when the content has been compressed.
	Size int `json:"size"`

	// Number of bytes saved.
	// Leave out this field if the information is not available.
	Compression int `json:"compression"`

	// MIME type of the response text
	// (value of the Content-Type response header).
	// The charset attribute of the MIME type is included (if available).
	MimeType string `json:"mimeType"`

	// Response body sent from the server or loaded from the browser cache.
	// This field is populated with textual content only.
	// The text field is either HTTP decoded text or a encoded
	// (e.g. "base64") representation of the response body.
	// Leave out this field if the information is not available.
	Text string `json:"text"`

	// Encoding used for response text field e.g "base64".
	// Leave out this field if the text field is HTTP decoded
	// (decompressed & unchunked), then trans-coded from its original
	// character set into UTF-8.
	Encoding string `json:"encoding,omitempty"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// This objects contains info about a request coming from browser cache.
type Cache struct {
	// BeforeRequest is the state of a cache entry before the request.
	// Leave out this field if the information is not available.
	BeforeRequest CacheEntryState `json:"beforeRequest,omitempty"`

	// AfterRequest is the state of a cache entry after the request.
	// Leave out this field if the information is not available.
	AfterRequest CacheEntryState `json:"afterRequest,omitempty"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// CacheEntryState contains information about a cache entry.
type CacheEntryState struct {
	// Expires is the expiration time of the cache entry.
	Expires string `json:"expires"`

	// LastAccess is the last time the cache entry was opened.
	LastAccess string `json:"lastAccess"`

	// Etag
	ETag string `json:"eTag"`

	// HitCount is the  number of times the cache entry has been opened.
	HitCount int `json:"hitCount"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}

// Timings describes various phases within request-response round trip.
// All times are specified in milliseconds.
// Note:
// The send, wait and receive timings are not optional and must have
// non-negative values.
type Timings struct {
	// Blocked is the time spent in a queue waiting for a network connection.
	// Use -1 if the timing does not apply to the current request.
	Blocked int `json:"blocked"`

	// DNS is the DNS resolution time. The time required to resolve a host name.
	// Use -1 if the timing does not apply to the current request.
	DNS int `json:"dns"`

	// Connect is the time required to create the TCP connection.
	// Use -1 if the timing does not apply to the current request.
	Connect int `json:"connect"`

	// Send is the time required to send HTTP request to the server.
	Send int `json:"send"`

	// Wait is the time spent waiting for a response from the server.
	Wait int `json:"wait"`

	// Receive is the time required to read the entire response from the server
	// (or cache).
	Receive int `json:"receive"`

	// SSL is the time required for SSL/TLS negotiation.
	// If this field is defined then the time is also included in the connect
	// field (to ensure backward compatibility with HAR 1.1).
	// Use -1 if the timing does not apply to the current request.
	SSL int `json:"ssl"`

	// Comment is a comment provided by the user or the application.
	Comment string `json:"comment,omitempty"`
}
