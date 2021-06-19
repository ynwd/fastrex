package fastrex

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpHandler_validate(t *testing.T) {
	type fields struct {
		routes      map[string]appRoute
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

func TestHandlerFunc_ServeHTTP(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		f    HandlerFunc
		args args
	}{
		{
			name: "success",
			f:    func(r1 Request, r2 Response) {},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}
