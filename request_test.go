package fastrex

import (
	"net/http"
	"net/http/httptest"
	"reflect"
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
			want: params,
		},
		{
			name:     "fail",
			args:     args{},
			incoming: httptest.NewRequest("GET", "/address/jakarta/cirebon", nil),
			routes: map[string]appRoute{
				"GET:/": {path: "/user/:name", method: "GET", handler: nil},
			},
			want: params,
		},
		{
			name: "fail",
			args: args{
				name: []string{"agus", "budi"},
			},
			incoming: httptest.NewRequest("GET", "/address/jakarta/cirebon", nil),
			routes:   map[string]appRoute{},
			want:     params,
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
