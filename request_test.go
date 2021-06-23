package fastrex

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRequest_Params(t *testing.T) {
	params := []string{}

	type args struct {
		name []string
	}
	tests := []struct {
		name     string
		args     args
		incoming *http.Request
		routes   map[string]appRoute
		want     []string
	}{
		{
			name:     "success",
			args:     args{},
			incoming: httptest.NewRequest("GET", "/user/agus/jakarta", nil),
			routes: map[string]appRoute{
				"GET:/": {path: "/user/:name/:address", method: "GET", handler: nil},
			},
			want: append(params, "agus", "jakarta"),
		},
		{
			name: "success",
			args: args{
				name: []string{"name"},
			},
			incoming: httptest.NewRequest("GET", "/user/agus/jakarta", nil),
			routes: map[string]appRoute{
				"GET:/": {path: "/user/:name/:address", method: "GET", handler: nil},
			},
			want: append(params, "agus"),
		},
		{
			name:     "fail",
			args:     args{},
			incoming: httptest.NewRequest("GET", "/address/jakarta", nil),
			routes: map[string]appRoute{
				"GET:/": {path: "/user/:name", method: "GET", handler: nil},
			},
			want: []string{},
		},
		{
			name:     "fail",
			args:     args{},
			incoming: httptest.NewRequest("GET", "/address/jakarta/cirebon", nil),
			routes: map[string]appRoute{
				"GET:/": {path: "/user/:name", method: "GET", handler: nil},
			},
			want: []string{},
		},
		{
			name: "fail",
			args: args{
				name: []string{"agus", "budi"},
			},
			incoming: httptest.NewRequest("GET", "/address/jakarta/cirebon", nil),
			routes:   map[string]appRoute{},
			want:     []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newRequest(tt.incoming, tt.routes)
			if got := h.Params(tt.args.name...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.Params() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Cookie(t *testing.T) {
	incomingCookie := http.Cookie{}
	incomingCookie.Name = "name"
	incomingCookie.Value = "agus"
	incoming := httptest.NewRequest("GET", "/", nil)
	incoming.AddCookie(&incomingCookie)

	expectedCookie := Cookie{}
	expectedCookie.Name("name").Value("agus")

	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		incoming *http.Request
		routes   map[string]appRoute
		want     Cookie
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				name: "name",
			},
			incoming: incoming,
			routes:   map[string]appRoute{},
			want:     expectedCookie,
			wantErr:  false,
		},
		{
			name: "failed",
			args: args{
				name: "agus",
			},
			incoming: incoming,
			routes:   map[string]appRoute{},
			want:     Cookie{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newRequest(tt.incoming, tt.routes)
			got, err := h.Cookie(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.Cookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.Cookie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Cookies(t *testing.T) {
	incomingCookie := http.Cookie{}
	incomingCookie.Name = "name"
	incomingCookie.Value = "agus"
	incoming := httptest.NewRequest("GET", "/", nil)
	incoming.AddCookie(&incomingCookie)

	expectedCookie := Cookie{}
	expectedCookie.Name("name").Value("agus")
	cookies := []Cookie{}
	cookies = append(cookies, expectedCookie)

	tests := []struct {
		name     string
		incoming *http.Request
		routes   map[string]appRoute
		want     []Cookie
	}{
		{
			name:     "success",
			incoming: incoming,
			routes:   map[string]appRoute{},
			want:     cookies,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newRequest(tt.incoming, tt.routes)
			if got := h.Cookies(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.Cookies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Context(t *testing.T) {
	t.Run("context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)

		var handler HandlerFunc = func(r1 Request, r2 Response) {
			r1.Context()
		}
		handler.ServeHTTP(nil, req, nil, nil)
		ctx := req.Context()
		wantCtx := context.Background()
		if ctx != wantCtx {
			t.Errorf("Request.Context() = %v, want %v", ctx, wantCtx)
		}
	})
}

func TestRequest_AddCookie(t *testing.T) {
	t.Run("Add cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)

		var handler HandlerFunc = func(r1 Request, r2 Response) {
			c := Cookie{}
			c.Name("name").Value("agus")
			r1.AddCookie(c)
		}
		handler.ServeHTTP(nil, req, nil, nil)

		got := Cookie{}
		c, _ := req.Cookie("name")
		got.Name(c.Name).Value(c.Value)

		want := Cookie{}
		want.Name("name").Value("agus")

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Request.AddCookie() = %v, want %v", got, want)
		}
	})
}

func TestRequest_BasicAuth(t *testing.T) {
	t.Run("basic auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.SetBasicAuth("agus", "password")
		ok := false
		var username, password string
		var handler HandlerFunc = func(r1 Request, r2 Response) {
			username, password, ok = r1.BasicAuth()
		}
		handler.ServeHTTP(nil, req, nil, nil)
		wantUsername := "agus"
		wantPassword := "password"
		wantOk := true
		if ok != wantOk && wantUsername != username && wantPassword != password {
			t.Errorf("Request.BasicAuth() = %v, want %v", ok, wantOk)
		}
	})
}

func TestRequest_Clone(t *testing.T) {
	t.Run("clone", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		var got Request
		var handler HandlerFunc = func(req Request, r2 Response) {
			ctx := context.Background()
			got = req.Clone(ctx)
		}

		handler.ServeHTTP(nil, req, nil, nil)

		want := *newRequest(req, nil)
		want = want.WithContext(req.Context())

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Request.Clone()\n = %v\n, want\n = %v\n", got, want)
		}
	})
}

func TestRequest_UserAgent(t *testing.T) {
	t.Run("user agent", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("User-Agent", "test")
		var got string

		var handler HandlerFunc = func(r1 Request, r2 Response) {
			got = r1.UserAgent()
		}
		handler.ServeHTTP(nil, req, nil, nil)

		want := "test"
		if got != want {
			t.Errorf("Request.Context() = %v, want %v", got, want)
		}
	})
}

func TestRequest_Referer(t *testing.T) {
	t.Run("referer", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Referer", "test")
		var got string

		var handler HandlerFunc = func(r1 Request, r2 Response) {
			got = r1.Referer()
		}
		handler.ServeHTTP(nil, req, nil, nil)

		want := "test"
		if got != want {
			t.Errorf("Request.Context() = %v, want %v", got, want)
		}
	})
}

func TestRequest_FormFile(t *testing.T) {
	t.Run("FormFile", func(t *testing.T) {
		path := "README.md" //The path to upload the file
		file, err := os.Open(path)
		if err != nil {
			t.Error(err)
		}

		defer file.Close()
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("my_file", filepath.Base(path))

		if err != nil {
			writer.Close()
			t.Error(err)
		}
		io.Copy(part, file)
		writer.Close()

		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		var handler HandlerFunc = func(req Request, r2 Response) {
			_, _, err = req.FormFile("my_file")
		}

		handler.ServeHTTP(res, req, nil, nil)

		if res.Code != http.StatusOK {
			t.Error("not 200")
		}

		if err != nil {
			t.Error(err.Error())
		}
	})
}

func TestRequest_Write(t *testing.T) {
	t.Run("write", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		var err error
		var braw bytes.Buffer
		var handler HandlerFunc = func(r1 Request, r2 Response) {
			err = r1.Write(&braw)
		}
		handler.ServeHTTP(nil, req, nil, nil)
		if err != nil {
			t.Errorf("Request.Write() = %v", err)
		}
	})
}

func TestRequest_WriteProxy(t *testing.T) {
	t.Run("write proxy", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		var got error
		var braw bytes.Buffer
		var handler HandlerFunc = func(r1 Request, r2 Response) {
			got = r1.WriteProxy(&braw)
		}
		handler.ServeHTTP(nil, req, nil, nil)
		if got != nil {
			t.Errorf("Request.WriteProxy() = %v", got)
		}
	})
}

func TestRequest_ErrorMiddleware(t *testing.T) {
	t.Run("error middleware", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		var r Request
		var handler HandlerFunc = func(r1 Request, r2 Response) {
			r = r1.ErrorMiddleware(errors.New("not found"), 404)
		}
		handler.ServeHTTP(nil, req, nil, nil)
		got := r.Context().Value(errMiddlewareKey)
		want := ErrMiddleware{
			Error: errors.New("not found"),
			Code:  404,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Request.ErrorMiddleware() = %v, want %v", got, want)
		}
	})
}

func TestRequest_MultipartReader(t *testing.T) {
	t.Run("MultipartReader", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", io.NopCloser(new(bytes.Buffer)))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=foo123")
		var err error
		var handler HandlerFunc = func(r1 Request, r2 Response) {
			_, err = r1.MultipartReader()
		}

		handler.ServeHTTP(nil, req, nil, nil)
		if err != nil {
			t.Errorf("Request.ErrorMiddleware() = %v", err)
		}
	})
}

func TestRequest_FormValue(t *testing.T) {
	t.Run("FormValue", func(t *testing.T) {
		req := httptest.NewRequest("POST", "http://www.google.com/", strings.NewReader("z=post"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		var got string
		var handler HandlerFunc = func(r1 Request, r2 Response) {
			got = r1.FormValue("z")
		}

		handler.ServeHTTP(nil, req, nil, nil)
		if got != "post" {
			t.Errorf("Request.FormValue() = %v", got)
		}
	})
}

func TestRequest_ParseFrom(t *testing.T) {
	t.Run("ParseFrom", func(t *testing.T) {
		req := httptest.NewRequest("POST", "http://www.google.com/", strings.NewReader("z=post"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		var err error
		var handler HandlerFunc = func(r1 Request, r2 Response) {
			err = r1.ParseForm()
		}

		handler.ServeHTTP(nil, req, nil, nil)
		if err != nil {
			t.Errorf("Request.ParseForm() = %v", err)
		}
	})
}

func TestRequest_ParseMultipartForm(t *testing.T) {
	t.Run("ParseMultipartForm", func(t *testing.T) {
		postData :=
			`--xxx
Content-Disposition: form-data; name="field1"

value1
--xxx
Content-Disposition: form-data; name="field2"

value2
--xxx
Content-Disposition: form-data; name="file"; filename="file"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data
--xxx--
`
		req := httptest.NewRequest("POST", "/", io.NopCloser(strings.NewReader(postData)))
		req.Header.Set("Content-Type", `multipart/form-data; boundary=xxx`)

		initialFormItems := map[string]string{
			"language": "Go",
			"name":     "gopher",
			"skill":    "go-ing",
			"field2":   "initial-value2",
		}

		req.Form = make(url.Values)
		for k, v := range initialFormItems {
			req.Form.Add(k, v)
		}

		var err error
		var handler HandlerFunc = func(r1 Request, r2 Response) {
			err = r1.ParseMultipartForm(25)
		}

		handler.ServeHTTP(nil, req, nil, nil)
		if err != nil {
			t.Errorf("Request.ParseMultipartForm() = %v", err)
		}
	})
}

func TestRequest_ProtoAtLeast(t *testing.T) {
	t.Run("ProtoAtLeast", func(t *testing.T) {
		req := httptest.NewRequest("POST", "http://www.google.com/", nil)
		var got bool
		var handler HandlerFunc = func(r1 Request, r2 Response) {
			got = r1.ProtoAtLeast(1, 1)
		}

		handler.ServeHTTP(nil, req, nil, nil)
		if !got {
			t.Errorf("Request.ProtoAtLeast() = %v", got)
		}
	})
}
