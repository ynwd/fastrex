package fastrex

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpHandler_validate(t *testing.T) {
	type fields struct {
		routes      map[string]AppRoute
		middlewares []Middleware
		logger      *log.Logger
		ctx         context.Context
	}
	type args struct {
		path     string
		incoming string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "success",
			args: args{
				path:     "/view/user/:id/view/:name",
				incoming: "/view/user/6/view/agus",
			},
			want: true,
		},
		{
			name: "success",
			args: args{
				path:     "/view/user/:id/view",
				incoming: "/view/user/6/view",
			},
			want: true,
		},
		{
			name: "success",
			args: args{
				path:     "/:id",
				incoming: "/x",
			},
			want: true,
		},
		{
			name: "fail",
			args: args{
				path:     "/view/:id",
				incoming: "/x",
			},
			want: false,
		},
		{
			name: "fail",
			args: args{
				path:     "/view/",
				incoming: "/x",
			},
			want: false,
		},
		{
			name: "fail",
			args: args{
				path:     "/",
				incoming: "/x",
			},
			want: false,
		},
		{
			name: "success",
			args: args{
				path:     "/user/:id([0-9]+)",
				incoming: "/user/9",
			},
			want: true,
		},
		{
			name: "fail",
			args: args{
				path:     "/user/:id([0-9]+)",
				incoming: "/user/agus",
			},
			want: false,
		},
		{
			name: "success",
			args: args{
				path:     "/user/:id()",
				incoming: "/user/agus",
			},
			want: true,
		},
		{
			name: "fail",
			args: args{
				path:     "/siap/:name/address/field",
				incoming: "/oke/agus/address/field",
			},
			want: false,
		},
		{
			name: "fail",
			args: args{
				path:     "/siap/:name/oke/oke",
				incoming: "/siap/agus/address/field",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &httpHandler{
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				logger:      tt.fields.logger,
				ctx:         tt.fields.ctx,
			}
			if got := h.validate(tt.args.path, tt.args.incoming); got != tt.want {
				t.Errorf("HttpHandler.validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpRouter_ServeHTTP(t *testing.T) {
	files := []string{}
	tmpl, _ := template.ParseFiles("template/index.html")
	nilTmpl, _ := template.ParseFiles("template/app.html")
	errorMiddleware := []Middleware{}
	errorMiddleware = append(errorMiddleware, func(r1 Request, r2 Response, n Next) {
		r1 = r1.ErrorMiddleware(http.ErrAbortHandler, http.StatusBadGateway)
		n(r1, r2)
	})

	normalMiddleware := []Middleware{}
	normalMiddleware = append(normalMiddleware, func(r1 Request, r2 Response, n Next) {
		n(r1, r2)
	})

	type fields struct {
		routes       map[string]AppRoute
		middlewares  []Middleware
		logger       *log.Logger
		ctx          context.Context
		filenames    []string
		template     *template.Template
		staticPath   string
		staticFolder string
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
				routes: map[string]AppRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {}},
				},
				logger: log.Default(),
			},
			args: args{
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]AppRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {}},
				},
				template: tmpl,
			},
			args: args{
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]AppRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {}},
				},
				template: nilTmpl,
			},
			args: args{
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]AppRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {}},
				},
				filenames: files,
			},
			args: args{
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				staticFolder: "static",
				staticPath:   "/",
				routes: map[string]AppRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/ok", nil),
			},
		},
		{
			name: "failed",
			fields: fields{
				routes: map[string]AppRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
					}},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/ok", nil),
			},
		},
		{
			name: "success - route middleware",
			fields: fields{
				routes: map[string]AppRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Json(7)
					}, middlewares: normalMiddleware},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success - app middleware",
			fields: fields{
				routes: map[string]AppRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Send("")
					}},
				},
				middlewares: normalMiddleware,
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest("GET", "/", nil),
			},
		},
		{
			name: "success",
			fields: fields{
				routes: map[string]AppRoute{
					"GET:/": {path: "/", method: "GET", handler: func(r1 Request, r2 Response) {
						r2.Send("")
					}},
				},
				middlewares: errorMiddleware,
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
				staticFolder: tt.fields.staticFolder,
				staticPath:   tt.fields.staticPath,
				routes:       tt.fields.routes,
				middlewares:  tt.fields.middlewares,
				logger:       tt.fields.logger,
				ctx:          tt.fields.ctx,
				filename:     tt.fields.filenames,
				template:     tt.fields.template,
			}
			r.ServeHTTP(tt.args.res, tt.args.req)
		})
	}
}
