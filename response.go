package fastrex

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

const (
	HeaderContentType   = "Content-Type"
	MimeApplicationJson = "application/json"
	MimeTextHtml        = "text/html; charset=UTF-8"
	serverlessFolder    = "serverless_function_source_code/"
)

type Response interface {
	// Get header
	Header() http.Header
	// Write writes the data to the connection as part of an HTTP reply.
	Write([]byte) (int, error)
	// WriteHeader sends an HTTP response header with the provided status code.
	WriteHeader(statusCode int)
	// Sets the response’s HTTP header field to value
	Set(string, string) Response
	// Sets the HTTP status for the response
	Status(int) Response
	// Sets the Content-Type HTTP header to the MIME type as determined by the specified type.
	Type(string) Response
	// Sends a JSON response.
	Json(data interface{})
	// Sends the HTTP response.
	Send(data interface{})
	// Sets a created cookie.
	Cookie(Cookie) Response
	// Clears the cookie specified by name.
	ClearCookie(cookie Cookie) Response
	// Redirects to the URL derived from the specified path, with specified status.
	Redirect(path string, code int)
	// Sets the response Location HTTP header to the specified path parameter.
	Location(path string) Response
	// Appends the specified value to the HTTP response header field.
	Append(key string, val string) Response
	// Renders a view and sends the rendered HTML string to the client.
	Render(args ...interface{}) error
	// SendFile()
	// SendStatus()
	// Jsonp()
	// Vary() Response
	// Download()
	// End()
	// Format() Response
	// Get() Response
	// Links() Response
	// Attachment() Response
}

func newResponse(w http.ResponseWriter, r *http.Request, t *template.Template) Response {
	return &httpResponse{
		w: w,
		r: r,
		s: http.StatusOK,
		t: t,
	}
}

type httpResponse struct {
	s int
	c *Cookie
	w http.ResponseWriter
	r *http.Request
	t *template.Template
}

func (h *httpResponse) Header() http.Header {
	return h.w.Header()
}

func (h *httpResponse) Write(data []byte) (int, error) {
	return h.w.Write(data)
}

func (h *httpResponse) WriteHeader(statusCode int) {
	h.w.WriteHeader(statusCode)
}

func (h *httpResponse) Type(httpType string) Response {
	h.w.Header().Set(HeaderContentType, httpType)
	return h
}

func (h *httpResponse) Set(field string, value string) Response {
	h.w.Header().Set(field, value)
	return h
}

func (h *httpResponse) Status(code int) Response {
	h.s = code
	return h
}

func (h *httpResponse) Send(data interface{}) {
	if h.c != nil {
		c := h.c.cookie()
		http.SetCookie(h.w, c)
	}
	h.w.WriteHeader(h.s)
	h.w.Write([]byte(data.(string)))
}

func (h *httpResponse) Json(data interface{}) {
	var jsonStr string

	switch data.(type) {
	case string, bool, int, int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128:
		jsonStr = fmt.Sprintf("%v", data)
	default:
		jsonStr = processStruct(data)
	}

	h.w.Header().Set(HeaderContentType, MimeApplicationJson)
	h.w.Write([]byte(jsonStr))
}

func processStruct(data interface{}) string {
	jsonByte, err := json.Marshal(data)
	if err != nil {
		return err.Error()
	}
	return string(jsonByte)
}

func (h *httpResponse) Cookie(cookie Cookie) Response {
	h.c = &cookie
	return h
}

func (h *httpResponse) ClearCookie(cookie Cookie) Response {
	c := Cookie{}
	if cookie.c.Name != "" {
		c.c = cookie.c
		c.MaxAge(-1)
		h.c = &c
	}
	return h
}

func (h *httpResponse) Redirect(url string, code int) {
	http.Redirect(h.w, h.r, url, code)
}

func (h *httpResponse) Location(path string) Response {
	h.w.Header().Set("Location", path)
	return h
}

func (h *httpResponse) Append(field string, value string) Response {
	h.w.Header().Add(field, value)
	return h
}

func (h *httpResponse) Render(args ...interface{}) error {
	if h.t == nil {
		return errors.New("empty template")
	}
	length := len(args)
	h.w.Header().Set(HeaderContentType, MimeTextHtml)
	if length == 0 || length == 1 {
		if length == 0 {
			return h.t.Execute(h.w, nil)
		}
		return h.t.Execute(h.w, args[0])
	} else if length == 2 {
		name := args[0].(string)
		data := args[1]
		if name == "" {
			return errors.New("empty template name")
		}
		return h.t.ExecuteTemplate(h.w, name, data)
	}
	return errors.New("invalid args")
}

// TODO:
// func (h *httpResponse) Attachment() Response {
// 	return h
// }

// func (h *httpResponse) Download() {
// }

// func (h *httpResponse) End() {
// }

// func (h *httpResponse) Format() Response {
// 	return h
// }

// func (h *httpResponse) Get() Response {
// 	return h
// }

// func (h *httpResponse) Jsonp() {
// }

// func (h *httpResponse) Links() Response {
// 	return h
// }

// func (h *httpResponse) SendFile() {
// }

// func (h *httpResponse) SendStatus() {
// }

// func (h *httpResponse) Vary() Response {
// 	return h
// }
