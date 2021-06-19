package fastrex

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_httpRouter_ServeHTTP(t *testing.T) {
	files := make([]string, 0)
	tmpl, _ := template.ParseFiles("template/index.html")
	nilTmpl, _ := template.ParseFiles("template/app.html")
	m := []Middleware{}
	m = append(m, func(r1 Request, r2 Response, n Next) {
		r1 = r1.ErrorMiddleware(http.ErrAbortHandler, http.StatusBadGateway)
		n(r1, r2)
	})

	m2 := []Middleware{}
	m2 = append(m2, func(r1 Request, r2 Response, n Next) {
		n(r1, r2)
	})

	type fields struct {
		routes      map[string]appRoute
		middlewares []Middleware
		logger      *log.Logger
		ctx         context.Context
		filenames   []string
		template    *template.Template
	}
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				ctx: context.Background(),
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r1.Context()
						c := Cookie{}
						c.Name("name").Value("agus").Path("/")
						r2.ClearCookie(c)
						r2.Cookie(c).
							Type(MimeApplicationJson).
							Set("x", "x").
							Status(200).
							Send("oke")
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				ctx: context.Background(),
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Redirect("/oke", 200)
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				ctx: context.Background(),
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Location("/oke")
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				ctx: context.Background(),
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Append("x", "y")
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				template: tmpl,
				ctx:      context.Background(),
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Render()
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				template: tmpl,
				ctx:      context.Background(),
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						data := map[string]interface{}{
							"title": "Cloud function app",
							"name":  "Agus",
						}
						r2.Render(data)
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				template: tmpl,
				ctx:      context.Background(),
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						data := map[string]interface{}{
							"title": "Cloud function app",
							"name":  "Agus",
						}
						r2.Render("n", data)
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				template: nilTmpl,
				ctx:      context.Background(),
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Render()
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				template: tmpl,
				ctx:      context.Background(),
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Render("", "")
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "invalid args of render",
			fields: fields{
				template: tmpl,
				ctx:      context.Background(),
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Render("", "", "")
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						data := map[string]interface{}{
							"title": "Learning Golang Web",
							"name":  "Batman",
						}
						r2.Json(data)
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						value := make(chan int)
						r2.Status(200).
							Json(value)
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "empty cookie",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r1.Cookies()
						r1.Cookie("cookie")
						r2.Status(200).
							Json("ok")
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r1.Header.Set("Cookie", "cookie=xxxx")
						r1.Cookies()
						r1.Cookie("cookie")
						r1.FormFile("x")
						r1.FormValue("x")
						r1.ProtoAtLeast(0, 1)
						r1.ParseForm()
						r1.ParseMultipartForm(1)
						r1.MultipartReader()
						r1.Referer()
						r1.UserAgent()
						r1.WithContext(context.Background())
						r1.AddCookie(Cookie{})
						r1.BasicAuth()
						r1.Clone(context.Background())
						r1.Write(r2)
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r1.WriteProxy(r2)
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						c := Cookie{}
						expiration := time.Now().Add(365 * 24 * time.Hour)
						c.Name("user").Path("/").Expires(expiration).MaxAge(10).Raw("ok").HttpOnly(true).Secure(true)
						c.Unparsed()
						c.SameSite(http.SameSiteDefaultMode)
						c.RawExpires()
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Header().Set("x", "y")
						r2.WriteHeader(200)
						r2.Write([]byte("ok"))
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Json(7)
					}, middlewares: m2},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/ok/ok", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Send("")
					}},
				},
				middlewares: m2,
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Send("")
					}},
				},
				logger:      log.Default(),
				ctx:         context.Background(),
				middlewares: m,
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Send("")
					}},
				},
				logger:      log.Default(),
				ctx:         context.Background(),
				middlewares: m,
				filenames:   append(files, "ok"),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]appRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Status(200).
							Send("ok")
					}},
				},
				logger:      log.Default(),
				ctx:         context.Background(),
				middlewares: m,
				filenames:   append(files, "app.html"),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &app{
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				logger:      tt.fields.logger,
				ctx:         tt.fields.ctx,
				filename:    tt.fields.filenames,
				template:    tt.fields.template,
			}
			r.ServeHTTP(tt.args.res, tt.args.req)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want App
	}{
		{
			name: "success",
			want: &app{
				routes:      map[string]appRoute{},
				middlewares: []Middleware{},
				server:      &http.Server{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpRouter_Get(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		r := &app{
			routes: map[string]appRoute{},
		}
		want := &app{
			routes: map[string]appRoute{
				"GET:/": {path: "/", method: "GET"},
			},
		}
		if got := r.Get("/", nil); !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.Get() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_httpRouter_Connect(t *testing.T) {
	t.Run("CONNECT", func(t *testing.T) {
		r := &app{
			routes: map[string]appRoute{},
		}
		want := &app{
			routes: map[string]appRoute{
				"CONNECT:/": {path: "/", method: "CONNECT"},
			},
		}
		got := r.Connect("/", nil)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.Connect() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_httpRouter_Delete(t *testing.T) {
	t.Run("DELETE", func(t *testing.T) {
		r := &app{
			routes: map[string]appRoute{},
		}
		want := &app{
			routes: map[string]appRoute{
				"DELETE:/": {path: "/", method: "DELETE"},
			},
		}
		got := r.Delete("/", nil)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.Delete() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_httpRouter_Head(t *testing.T) {
	t.Run("HEAD", func(t *testing.T) {
		r := &app{
			routes: map[string]appRoute{},
		}
		want := &app{
			routes: map[string]appRoute{
				"HEAD:/": {path: "/", method: "HEAD"},
			},
		}
		got := r.Head("/", nil)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.Head() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_httpRouter_Put(t *testing.T) {
	t.Run("PUT", func(t *testing.T) {
		r := &app{
			routes: map[string]appRoute{},
		}
		want := &app{
			routes: map[string]appRoute{
				"PUT:/": {path: "/", method: "PUT"},
			},
		}
		got := r.Put("/", nil)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.Put() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_httpRouter_Patch(t *testing.T) {
	t.Run("PATCH", func(t *testing.T) {
		r := &app{
			routes: map[string]appRoute{},
		}
		want := &app{
			routes: map[string]appRoute{
				"PATCH:/": {path: "/", method: "PATCH"},
			},
		}
		got := r.Patch("/", nil)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.Patch() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_httpRouter_Trace(t *testing.T) {
	t.Run("TRACE", func(t *testing.T) {
		r := &app{
			routes: map[string]appRoute{},
		}
		want := &app{
			routes: map[string]appRoute{
				"TRACE:/": {path: "/", method: "TRACE"},
			},
		}
		got := r.Trace("/", nil)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.Trace() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_httpRouter_Post(t *testing.T) {
	t.Run("POST", func(t *testing.T) {
		r := &app{
			routes: map[string]appRoute{},
		}
		want := &app{
			routes: map[string]appRoute{
				"POST:/": {path: "/", method: "POST"},
			},
		}
		got := r.Post("/", nil)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.Post() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_httpRouter_Method_withHandlerAndMiddleware(t *testing.T) {
	handler := func(req Request, res Response) {}
	middleware := func(req Request, res Response, next Next) {
		next(req, res)
	}
	middlewares := []Middleware{}
	middlewares = append(middlewares, middleware)

	t.Run("GET", func(t *testing.T) {
		r := &app{
			routes: map[string]appRoute{},
		}
		want := &app{
			routes: map[string]appRoute{
				"GET:/": {path: "/", method: "GET", handler: handler, middlewares: middlewares},
			},
		}
		if r.Get("/", handler, middleware); !checkLen(r.routes, want.routes) {
			t.Errorf("httpRouter.Get() = %v, want %v", r.routes, want.routes)
		}
	})
}

func checkLen(got map[string]appRoute, want map[string]appRoute) bool {
	return len(got) == len(want)
}

func Test_httpRouter_Log(t *testing.T) {
	type fields struct {
		routes      map[string]appRoute
		middlewares []Middleware
		logger      *log.Logger
		ctx         context.Context
	}
	type args struct {
		logger *log.Logger
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   App
	}{
		{
			name:   "success",
			fields: fields{},
			args:   args{},
			want:   &app{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &app{
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				logger:      tt.fields.logger,
				ctx:         tt.fields.ctx,
			}
			if got := r.Log(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpRouter.Log() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpRouter_Ctx(t *testing.T) {
	type fields struct {
		routes      map[string]appRoute
		middlewares []Middleware
		logger      *log.Logger
		ctx         context.Context
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   App
	}{
		{
			name:   "success",
			fields: fields{},
			args:   args{},
			want:   &app{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &app{
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				logger:      tt.fields.logger,
				ctx:         tt.fields.ctx,
			}
			if got := r.Ctx(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpRouter.Ctx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpRouter_handler(t *testing.T) {
	type fields struct {
		routes      map[string]appRoute
		middlewares []Middleware
		logger      *log.Logger
		ctx         context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   http.Handler
	}{
		{
			name:   "success",
			fields: fields{},
			want:   &httpHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &app{
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				logger:      tt.fields.logger,
				ctx:         tt.fields.ctx,
			}
			if got := r.handler(false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpRouter.handler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpRouter_Listen_NonTLS(t *testing.T) {
	r := New()
	go func() {
		r.Listen(3000)
	}()
}

func Test_httpRouter_Listen_NonTLS_with_template(t *testing.T) {
	r := New()
	r.Template("template/index.html")
	go func() {
		r.Listen(3000)
	}()
}

func Test_httpRouter_Listen_NonTLS_with_err_template(t *testing.T) {
	r := New()
	r.Template("index.html")
	go func() {
		r.Listen(3000)
	}()
}

func Test_httpRouter_Listen_NonTLS_with_static(t *testing.T) {
	r := New()
	r.Static("static")
	go func() {
		r.Listen(3000)
	}()
}

func Test_httpRouter_Listen_NonTLS_with_err_static(t *testing.T) {
	r := New()
	r.Static("static", "/")
	go func() {
		r.Listen(3000)
	}()
}

func Test_httpRouter_Listen_NonTLS_same_port(t *testing.T) {
	r := New()
	r.Listen(3000)
	go func() {
		r.Listen(3000)
	}()
}

func Test_httpRouter_Listen_NonTLS_with_host(t *testing.T) {
	r := New()
	r.Host("localhost")
	go func() {
		r.Listen(3000)
	}()
}

func Test_httpRouter_Listen_NonTLS_withCallback(t *testing.T) {
	r := New()
	go func() {
		r.Listen(3000, func(err error) {
			if err != nil {
				fmt.Println(err.Error())
			}
		})
	}()
}

func Test_httpRouter_Close(t *testing.T) {
	r := New()
	wantErr := false
	if err := r.Close(); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_NonTLS_withErrorCallback(t *testing.T) {
	r := New()
	wantErr := true
	if err := r.Listen(3000, ""); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_TLS(t *testing.T) {
	r := New()
	go func() {
		r.Listen(3000, "", "")
	}()
}

func Test_httpRouter_Listen_TLS_with_host(t *testing.T) {
	r := New()
	r.Host("localhost")
	go func() {
		r.Listen(3000, "", "")
	}()
}

func Test_httpRouter_Listen_TLS_withCallback(t *testing.T) {
	r := New()
	go func() {
		r.Listen(3000, "", "", func(err error) {
			if err != nil {
				fmt.Println(err.Error())
			}
		})
	}()
}

func Test_httpRouter_Listen_TLS_withErrorKey(t *testing.T) {
	r := New()
	wantErr := true
	if err := r.Listen(3000, 0, ""); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_TLS_withErrorSecret(t *testing.T) {
	r := New()
	wantErr := true
	if err := r.Listen(3000, "", 0); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_TLS_withErrorCallback(t *testing.T) {
	r := New()
	wantErr := true
	if err := r.Listen(3000, "", "", ""); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_withInvalidArgs(t *testing.T) {
	r := New()
	wantErr := true
	if err := r.Listen(3000, "", "", "", ""); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_Shutdown(t *testing.T) {
	r := New()
	ctx, cancel := context.WithCancel(context.Background())
	r.RegisterOnShutdown(cancel)
	r.SetKeepAlivesEnabled(true)
	go func() {
		r.Listen(3000)
		defer r.Shutdown(ctx)
	}()
}

func Test_httpRouter_Use(t *testing.T) {
	type fields struct {
		routes      map[string]appRoute
		middlewares []Middleware
		logger      *log.Logger
		ctx         context.Context
		server      *http.Server
	}
	type args struct {
		m Middleware
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name:   "success",
			fields: fields{},
			args:   args{},
			want:   0,
		},
		{
			name:   "success",
			fields: fields{},
			args: args{
				m: func(r1 Request, r2 Response, n Next) {},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &app{
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				logger:      tt.fields.logger,
				ctx:         tt.fields.ctx,
				server:      tt.fields.server,
			}
			r.Use(tt.args.m)
			if got := len(r.middlewares); got != tt.want {
				t.Errorf("httpRouter.Use() = %v, want %v", got, tt.want)
			}
		})
	}
}
