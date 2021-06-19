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
	Method           string
	URL              *url.URL
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	Body             io.ReadCloser
	GetBody          func() (io.ReadCloser, error)
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Host             string
	Form             url.Values
	PostForm         url.Values
	MultipartForm    *multipart.Form
	Trailer          http.Header
	RemoteAddr       string
	RequestURI       string
	TLS              *tls.ConnectionState
	Cancel           <-chan struct{}
	Response         *http.Response
	ctx              context.Context
	r                *http.Request
	routes           map[string]appRoute
}

func newRequest(r *http.Request, routes map[string]appRoute) *Request {
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
		routes:           routes,
	}
}

func (h *Request) Context() context.Context {
	return h.ctx
}

func (h *Request) AddCookie(c Cookie) {
	h.r.AddCookie(&c.c)
}

func (h *Request) BasicAuth() (username string, password string, ok bool) {
	return h.r.BasicAuth()
}

func (h *Request) Clone(ctx context.Context) *http.Request {
	return h.r.Clone(ctx)
}

func (h *Request) Cookie(name string) (Cookie, error) {
	c := Cookie{}
	cookie, err := h.r.Cookie(name)
	if err != nil {
		return c, err
	}
	c.c = *cookie
	return c, nil
}

func (h *Request) Cookies() []Cookie {
	cookies := []Cookie{}

	for _, v := range h.r.Cookies() {
		c := Cookie{}
		c.c = *v
		cookies = append(cookies, c)
	}

	return cookies
}

func (h *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return h.r.FormFile(key)
}

func (h *Request) FormValue(key string) string {
	return h.r.FormValue(key)
}

func (h *Request) ProtoAtLeast(major int, minor int) bool {
	return h.r.ProtoAtLeast(major, minor)
}

func (h *Request) ParseForm() error {
	return h.r.ParseForm()
}

func (h *Request) ParseMultipartForm(maxMemory int64) error {
	return h.r.ParseMultipartForm(maxMemory)
}

func (h *Request) MultipartReader() (*multipart.Reader, error) {
	return h.r.MultipartReader()
}

func (h *Request) Write(w io.Writer) error {
	return h.r.Write(w)
}

func (h *Request) WriteProxy(w io.Writer) error {
	return h.r.WriteProxy(w)
}

func (h *Request) WithContext(ctx context.Context) Request {
	r := h.r.WithContext(ctx)
	return *newRequest(r, h.routes)
}

func (h *Request) UserAgent() string {
	return h.r.UserAgent()
}

func (h *Request) Referer() string {
	return h.r.Referer()
}

func (h *Request) ErrorMiddleware(e error, code int) Request {
	err := ErrMiddleware{
		Error: e,
		Code:  code,
	}
	r := h.r.WithContext(context.WithValue(h.ctx, errMiddlewareKey, err))
	return *newRequest(r, h.routes)
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

	for _, route := range h.routes {
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
