package fastrex

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"testing"
)

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
		r.Listen(3000, func(err error) {})
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

func Test_httpRouter_Listen_NonTLS_with_error_template(t *testing.T) {
	r := New()
	r.Template("index.htm")
	wantErr := true
	if err := r.Listen(3000); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_TLS(t *testing.T) {
	r := New()
	wantErr := true
	if err := r.Listen(3000, "certFile", "keyFile"); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_TLS_with_host(t *testing.T) {
	r := New()
	r.Host("localhost")
	wantErr := true
	if err := r.Listen(3000, "certFile", "keyFile"); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_TLS_withCallback(t *testing.T) {
	r := New()
	wantErr := true
	if err := r.Listen(3000, "certFile", "keyFile", func(err error) {}); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_TLS_withErrorKey(t *testing.T) {
	r := New()
	wantErr := true
	if err := r.Listen(3000, 0, "keyFile"); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_TLS_withErrorSecret(t *testing.T) {
	r := New()
	wantErr := true
	if err := r.Listen(3000, "certFile", 0); (err != nil) != wantErr {
		t.Errorf("httpRouter.Listen() error = %v, wantErr %v", err, wantErr)
	}
}

func Test_httpRouter_Listen_TLS_withErrorCallback(t *testing.T) {
	r := New()
	wantErr := true
	if err := r.Listen(3000, "certFile", "keyFile", ""); (err != nil) != wantErr {
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
