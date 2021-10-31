package fastrex

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_httpResponse_Send(t *testing.T) {
	tests := []struct {
		name       string
		handler    HandlerFunc
		wantBody   string
		wantStatus int
		wantHeader http.Header
	}{
		{
			name: "json",
			handler: func(r1 Request, r2 Response) {
				r2.Json("ping")
			},
			wantBody:   "ping",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "json with struct",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus")
				r2.Cookie(c)
				data := map[string]interface{}{
					"title": "Learning Golang Web",
					"name":  "Batman",
				}
				r2.Json(data)
			},
			wantBody:   `{"name":"Batman","title":"Learning Golang Web"}`,
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
				"Set-Cookie":   {"name=agus"},
			},
		},
		{
			name: "json with invalid tipe",
			handler: func(r1 Request, r2 Response) {
				value := make(chan int)
				r2.Json(value)
			},
			wantBody:   `json: unsupported type: chan int`,
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "send",
			handler: func(r1 Request, r2 Response) {
				r2.Send("ping")
			},
			wantBody:   "ping",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
			},
		},
		{
			name: "send with status",
			handler: func(r1 Request, r2 Response) {
				r2.Status(300).Send("ping")
			},
			wantBody:   "ping",
			wantStatus: 300,
			wantHeader: map[string][]string{},
		},
		{
			name: "send with header",
			handler: func(r1 Request, r2 Response) {
				r2.Set("x", "x").Send("ping")
			},
			wantBody:   "ping",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"X":            {"x"},
			},
		},
		{
			name: "send with header",
			handler: func(r1 Request, r2 Response) {
				r2.Header().Set("x", "x")
				r2.Send("ping")
			},
			wantBody:   "ping",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"X":            {"x"},
			},
		},
		{
			name: "send with type",
			handler: func(r1 Request, r2 Response) {
				r2.Type("application/json").Send("ping")
			},
			wantBody:   "ping",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "send with append",
			handler: func(r1 Request, r2 Response) {
				r2.Append("x", "y").Send("ping")
			},
			wantBody:   "ping",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"X":            {"y"},
			},
		},
		{
			name: "send with cookie",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus")
				r2.Cookie(c).Send("ping")
			},
			wantBody:   "ping",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"Set-Cookie":   {"name=agus"},
			},
		},
		{
			name: "send with clear cookie",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus")
				r2.ClearCookie(c).Send("ping")
			},
			wantBody:   "ping",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"Set-Cookie":   {"name=agus; Max-Age=0"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			tt.handler.ServeHTTP(w, req, nil, nil, nil, map[string]interface{}{})

			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)

			if got := string(body); !reflect.DeepEqual(got, tt.wantBody) {
				t.Errorf("Request.Params() = %v, want %v", got, tt.wantBody)
			}
			if status := resp.StatusCode; !reflect.DeepEqual(status, tt.wantStatus) {
				t.Errorf("Request.Params() = %v, want %v", status, tt.wantStatus)
			}
			if header := resp.Header; !reflect.DeepEqual(header, tt.wantHeader) {
				t.Errorf("Request.Params() = %v, want %v", header, tt.wantHeader)
			}
		})
	}
}

func Test_httpResponse_Write(t *testing.T) {
	tests := []struct {
		name       string
		handler    HandlerFunc
		wantBody   string
		wantStatus int
		wantHeader http.Header
	}{
		{
			name: "write",
			handler: func(r1 Request, r2 Response) {
				_, err := r2.Write([]byte("ping"))
				if err != nil {
					panic(err)
				}
			},
			wantBody:   "ping",
			wantStatus: 200,
			wantHeader: map[string][]string{},
		},
		{
			name: "write with writeheader",
			handler: func(r1 Request, r2 Response) {
				_, err := r2.WriteHeader(300).Write([]byte("ping"))
				if err != nil {
					panic(err)
				}
			},
			wantBody:   "ping",
			wantStatus: 300,
			wantHeader: map[string][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			tt.handler.ServeHTTP(w, req, nil, nil, nil, map[string]interface{}{})

			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)

			if got := string(body); !reflect.DeepEqual(got, tt.wantBody) {
				t.Errorf("Request.Params() = %v, want %v", got, tt.wantBody)
			}
			if status := resp.StatusCode; !reflect.DeepEqual(status, tt.wantStatus) {
				t.Errorf("Request.Params() = %v, want %v", status, tt.wantStatus)
			}
			if header := resp.Header; !reflect.DeepEqual(header, tt.wantHeader) {
				t.Errorf("Request.Params() = %v, want %v", header, tt.wantHeader)
			}
		})
	}
}

func Test_httpResponse_Location(t *testing.T) {
	tests := []struct {
		name       string
		handler    HandlerFunc
		wantHeader http.Header
	}{
		{
			name: "location",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus")
				r2.Cookie(c)
				r2.Location("/newlocation")
			},
			wantHeader: map[string][]string{
				"Location":   {"/newlocation"},
				"Set-Cookie": {"name=agus"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			w := httptest.NewRecorder()
			tt.handler.ServeHTTP(w, req, nil, nil, nil, map[string]interface{}{})
			resp := w.Result()
			if header := resp.Header; !reflect.DeepEqual(header, tt.wantHeader) {
				t.Errorf("Request.Params() = %v, want %v", header, tt.wantHeader)
			}
		})
	}
}

func Test_httpResponse_Redirect(t *testing.T) {
	tests := []struct {
		name       string
		handler    HandlerFunc
		wantStatus int
		wantHeader http.Header
	}{
		{
			name: "location",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus")
				r2.Cookie(c)
				r2.Redirect("/newlocation", 300)
			},
			wantHeader: map[string][]string{
				"Location":     {"/newlocation"},
				"Content-Type": {"text/html; charset=utf-8"},
				"Set-Cookie":   {"name=agus"},
			},
			wantStatus: 300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			w := httptest.NewRecorder()
			tt.handler.ServeHTTP(w, req, nil, nil, nil, map[string]interface{}{})
			resp := w.Result()
			if status := resp.StatusCode; !reflect.DeepEqual(status, tt.wantStatus) {
				t.Errorf("Request.Params() = %v, want %v", status, tt.wantStatus)
			}
			if header := resp.Header; !reflect.DeepEqual(header, tt.wantHeader) {
				t.Errorf("Request.Params() = %v, want %v", header, tt.wantHeader)
			}
		})
	}
}

func Test_httpResponse_Render(t *testing.T) {
	tests := []struct {
		name       string
		handler    HandlerFunc
		wantBody   string
		wantStatus int
		wantHeader http.Header
		template   string
		serverless bool
	}{
		{
			name:     "render serverless",
			template: "app.html",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus")
				r2.Cookie(c)
				err := r2.Render()
				if err != nil {
					panic(err)
				}
			},
			serverless: true,
			wantBody:   "<html></html>",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/html; charset=UTF-8"},
				"Set-Cookie":   {"name=agus"},
			},
		},
		{
			name:     "render serverless with data",
			template: "app.html",
			handler: func(r1 Request, r2 Response) {
				data := map[string]interface{}{
					"title": "app",
					"name":  "web",
				}
				err := r2.Render(data)
				if err != nil {
					panic(err)
				}
			},
			serverless: true,
			wantBody:   "<html>appweb</html>",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/html; charset=UTF-8"},
			},
		},
		{
			name:     "render serverless with name and data",
			template: "index.html",
			handler: func(r1 Request, r2 Response) {
				data := map[string]interface{}{
					"title": "app",
					"name":  "web",
				}
				err := r2.Render("index", data)
				if err != nil {
					panic(err)
				}
			},
			serverless: true,
			wantBody:   "<html>appweb</html>",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/html; charset=UTF-8"},
			},
		},
		{
			name:     "render serverless with empty name and data",
			template: "index.html",
			handler: func(r1 Request, r2 Response) {
				data := map[string]interface{}{
					"title": "app",
					"name":  "web",
				}
				r2.Render("", data)
			},
			serverless: true,
			wantBody:   "",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/html; charset=UTF-8"},
			},
		},
		{
			name:     "render serverless with invalid args",
			template: "index.html",
			handler: func(r1 Request, r2 Response) {
				data := map[string]interface{}{
					"title": "app",
					"name":  "web",
				}
				r2.Render("", data, "")
			},
			serverless: true,
			wantBody:   "",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/html; charset=UTF-8"},
			},
		},
		{
			name:     "render localhost",
			template: "template/app.html",
			handler: func(r1 Request, r2 Response) {
				r2.Render()
			},
			serverless: false,
			wantBody:   "<html></html>",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/html; charset=UTF-8"},
			},
		},
		{
			name:     "render localhost empty template",
			template: "",
			handler: func(r1 Request, r2 Response) {
				r2.Render()
			},
			serverless: false,
			wantBody:   "",
			wantStatus: 200,
			wantHeader: map[string][]string{},
		},
		{
			name:     "render localhost with data",
			template: "template/app.html",
			handler: func(r1 Request, r2 Response) {
				data := map[string]interface{}{
					"title": "app",
					"name":  "web",
				}
				r2.Render(data)
			},
			serverless: false,
			wantBody:   "<html>appweb</html>",
			wantStatus: 200,
			wantHeader: map[string][]string{
				"Content-Type": {"text/html; charset=UTF-8"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			app := New()
			app.Template(tt.template)
			app.Get("/", Handler(tt.handler))
			app.Serverless(tt.serverless)
			app.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			if got := string(body); !reflect.DeepEqual(got, tt.wantBody) {
				t.Errorf("Request.Params() = %v, want %v", got, tt.wantBody)
			}
			if status := resp.StatusCode; !reflect.DeepEqual(status, tt.wantStatus) {
				t.Errorf("Request.Params() = %v, want %v", status, tt.wantStatus)
			}
			if header := resp.Header; !reflect.DeepEqual(header, tt.wantHeader) {
				t.Errorf("Request.Params() = %v, want %v", header, tt.wantHeader)
			}
		})
	}
}
