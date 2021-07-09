package fastrex

import (
	"context"
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Cancel           <-chan struct{}
	GetBody          func() (io.ReadCloser, error)
	MultipartForm    *multipart.Form
	TLS              *tls.ConnectionState
	URL              *url.URL
	Response         *http.Response
	r                *http.Request
	Body             io.ReadCloser
	Trailer          http.Header
	Header           http.Header
	ctx              context.Context
	Form             url.Values
	PostForm         url.Values
	Routes           map[string]appRoute
	container        map[string]interface{}
	TransferEncoding []string
	Close            bool
	Serverless       bool
	RemoteAddr       string
	RequestURI       string
	Host             string
	Proto            string
	Method           string
	ProtoMajor       int
	ProtoMinor       int
	ContentLength    int64
}

// Context returns the request's context. To change the context, use WithContext.
//
// The returned context is always non-nil; it defaults to the background context.
//
// For outgoing client requests, the context controls cancellation.
//
// For incoming server requests, the context is canceled when the client's connection closes,
// the request is canceled (with HTTP/2), or when the ServeHTTP method returns.
func (h *Request) Context() context.Context {
	return h.r.Context()
}

// AddCookie adds a cookie to the request. Per RFC 6265 section 5.4,
// AddCookie does not attach more than one Cookie header field.
// That means all cookies, if any, are written into the same line, separated by semicolon.
// AddCookie only sanitizes c's name and value, and does not sanitize a Cookie header already
// present in the request.
func (h *Request) AddCookie(c Cookie) {
	h.r.AddCookie(&c.c)
}

// BasicAuth returns the username and password provided in the request's Authorization header,
// if the request uses HTTP Basic Authentication. See RFC 2617, Section 2.
func (h *Request) BasicAuth() (username string, password string, ok bool) {
	return h.r.BasicAuth()
}

// Clone returns a deep copy of r with its context changed to ctx.
// The provided ctx must be non-nil.
//
// For an outgoing client request, the context controls the entire
// lifetime of a request and its response: obtaining a connection,
// sending the request, and reading the response headers and body.
func (h *Request) Clone(ctx context.Context) Request {
	return *newRequest(h.r.Clone(ctx), h.Routes, h.Serverless, h.container)
}

// FormFile returns the first file for the provided form key.
// FormFile calls ParseMultipartForm and ParseForm if necessary.
func (h *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return h.r.FormFile(key)
}

// FormValue returns the first value for the named component of the query.
// POST and PUT body parameters take precedence over URL query string values.
// FormValue calls ParseMultipartForm and ParseForm if necessary and ignores
// any errors returned by these functions. If key is not present, FormValue
// returns the empty string. To access multiple values of the same key,
// call ParseForm and then inspect Request.Form directly.
func (h *Request) FormValue(key string) string {
	return h.r.FormValue(key)
}

// ProtoAtLeast reports whether the HTTP protocol used in the request is at least major.minor.
func (h *Request) ProtoAtLeast(major int, minor int) bool {
	return h.r.ProtoAtLeast(major, minor)
}

// ParseForm populates r.Form and r.PostForm.
//
// For all requests, ParseForm parses the raw
// query from the URL and updates r.Form.
//
// For POST, PUT, and PATCH requests, it also
// reads the request body, parses it as a form
// and puts the results into both r.PostForm and
// r.Form. Request body parameters take precedence
// over URL query string values in r.Form.
//
// If the request Body's size has not already been
// limited by MaxBytesReader, the size is capped at 10MB.
//
// For other HTTP methods, or when the Content-Type
// is not application/x-www-form-urlencoded, the request
// Body is not read, and r.PostForm is initialized to a
// non-nil, empty value.
//
// ParseMultipartForm calls ParseForm automatically.
// ParseForm is idempotent.
func (h *Request) ParseForm() error {
	return h.r.ParseForm()
}

// ParseMultipartForm parses a request body as multipart/form-data.
// The whole request body is parsed and up to a total of maxMemory
// bytes of its file parts are stored in memory, with the remainder
// stored on disk in temporary files. ParseMultipartForm calls
// ParseForm if necessary. After one call to ParseMultipartForm,
// subsequent calls have no effect.
func (h *Request) ParseMultipartForm(maxMemory int64) error {
	return h.r.ParseMultipartForm(maxMemory)
}

// MultipartReader returns a MIME multipart reader if this is a
// multipart/form-data or a multipart/mixed POST request,
// else returns nil and an error. Use this function instead of
// ParseMultipartForm to process the request body as a stream.
func (h *Request) MultipartReader() (*multipart.Reader, error) {
	return h.r.MultipartReader()
}

// Write writes an HTTP/1.1 request, which is the header and body,
// in wire format. This method consults the following fields of the request:
//  Host
//  URL
//  Method (defaults to "GET")
//  Header
//  ContentLength
//  TransferEncoding
//  Body
// If Body is present, Content-Length is <= 0 and TransferEncoding
// hasn't been set to "identity",
// Write adds "Transfer-Encoding: chunked" to the header.
// Body is closed after it is sent.
func (h *Request) Write(w io.Writer) error {
	return h.r.Write(w)
}

// WriteProxy is like Write but writes the request in
// the form expected by an HTTP proxy. In particular,
// WriteProxy writes the initial Request-URI line of
// the request with an absolute URI, per section 5.3
// of RFC 7230, including the scheme and host.
// In either case, WriteProxy also writes a Host header,
// using either r.Host or r.URL.Host.
func (h *Request) WriteProxy(w io.Writer) error {
	return h.r.WriteProxy(w)
}

// WithContext returns a shallow copy of r with its context
// changed to ctx. The provided ctx must be non-nil.
//
// For outgoing client request, the context controls
// the entire lifetime of a request and its response:
// obtaining a connection, sending the request, and reading
// the response headers and body.
//
// To create a new request with a context, use
// NewRequestWithContext. To change the context of a request,
// such as an incoming request you want to modify before
// sending back out, use Request.Clone. Between those two uses,
// it's rare to need WithContext.
func (h *Request) WithContext(ctx context.Context) Request {
	r := h.r.WithContext(ctx)
	return *newRequest(r, h.Routes, h.Serverless, h.container)
}

// UserAgent returns the client's User-Agent, if sent in the request.
func (h *Request) UserAgent() string {
	return h.r.UserAgent()
}

// Referer returns the referring URL, if sent in the request.
//
// Referer is misspelled as in the request itself, a mistake
// from the earliest days of HTTP. This value can also be fetched
// from the Header map as Header["Referer"]; the benefit of
// making it available as a method is that the compiler can
// diagnose programs that use the alternate (correct English)
// spelling req.Referer() but cannot diagnose programs that
// use Header["Referer"].
func (h *Request) Referer() string {
	return h.r.Referer()
}

func (h *Request) GetDependency(name string) interface{} {
	return h.container[name]
}

func (h *Request) SetDependency(name string, content interface{}) *Request {
	h.container[name] = content
	return h
}

func (h *Request) ErrorMiddleware(e error, code int) Request {
	err := ErrMiddleware{
		Error: e,
		Code:  code,
	}
	r := h.r.WithContext(context.WithValue(h.ctx, errMiddlewareKey, err))
	return *newRequest(r, h.Routes, h.Serverless, h.container)
}

func newRequest(r *http.Request, routes map[string]appRoute, serverless bool, container map[string]interface{}) *Request {
	return &Request{
		Method:           r.Method,
		URL:              r.URL,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Header:           r.Header,
		Body:             r.Body,
		GetBody:          r.GetBody,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Close:            r.Close,
		Host:             r.Host,
		Form:             r.Form,
		PostForm:         r.PostForm,
		MultipartForm:    r.MultipartForm,
		Trailer:          r.Header,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
		TLS:              r.TLS,
		Response:         r.Response,
		ctx:              r.Context(),
		r:                r,
		Routes:           routes,
		Serverless:       serverless,
		container:        container,
	}
}

// Cookie returns the named cookie provided in the request or ErrNoCookie if not found.
// If multiple cookies match the given name, only one cookie will be returned.
func (h *Request) Cookie(name string) (Cookie, error) {
	c := Cookie{}
	cookie, err := h.r.Cookie(name)
	if err != nil {
		return c, err
	}
	c.c = *cookie
	return c, nil
}

// Cookies parses and returns the HTTP cookies sent with the request.
func (h *Request) Cookies() []Cookie {
	cookies := []Cookie{}

	for _, v := range h.r.Cookies() {
		c := Cookie{}
		c.c = *v
		cookies = append(cookies, c)
	}

	return cookies
}

func (h *Request) Params(name ...string) []string {
	return h.getParams(name)
}

func (h *Request) getParams(name []string) []string {
	params := []string{}
	if len(name) > 1 {
		return params
	}

	incoming := split(h.r.URL.Path)
	length := len(incoming)

	for _, route := range h.Routes {
		routeChunks := split(route.path)
		routeLength := len(routeChunks)
		if length != routeLength {
			return []string{}
		}

		valid := parsePath(routeChunks, incoming)
		if !valid {
			return []string{}
		}

		if len(name) == 1 {
			params = append(params, getNamedParamItem(routeChunks, incoming, name[0])...)
		} else {
			params = append(params, getParamItem(routeChunks, incoming)...)
		}
	}
	return params
}

func getNamedParamItem(routeChunks []string, incoming []string, name string) []string {
	params := []string{}
	for idx, item := range routeChunks {
		if strings.Contains(item, splitter) && strings.Contains(item, splitter+name) {
			params = append(params, incoming[idx])
		}
	}
	return params
}

func getParamItem(routeChunks []string, incoming []string) []string {
	params := []string{}
	for idx, item := range routeChunks {
		if strings.Contains(item, splitter) {
			params = append(params, incoming[idx])
		}
	}
	return params
}
