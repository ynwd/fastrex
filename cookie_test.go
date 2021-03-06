package fastrex

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestCookie(t *testing.T) {
	expiration := time.Date(2011, 11, 23, 1, 5, 3, 0, time.UTC)
	expirationStr := expiration.Format("Mon, 2 Jan 2006 15:04:05 GMT")
	tests := []struct {
		name    string
		handler HandlerFunc
		want    http.Header
	}{
		{
			name: "domain",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Domain("localhost").Name("name").Value("agus")
				r2.Cookie(c)
				c.GetDomain()
				c.GetValue()
				r2.Send("")
			},
			want: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"Set-Cookie":   {"name=agus; Domain=localhost"},
			},
		},
		{
			name: "name",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus")
				r2.Cookie(c)
				c.GetName()
				r2.Send("")
			},
			want: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"Set-Cookie":   {"name=agus"},
			},
		},
		{
			name: "HttpOnly",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus").HttpOnly(true)
				r2.Cookie(c)
				c.GetHttpOnly()
				r2.Send("")
			},
			want: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"Set-Cookie":   {"name=agus; HttpOnly"},
			},
		},
		{
			name: "Path",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus").Path("/")
				r2.Cookie(c)
				c.GetPath()
				c.GetMaxAge()
				r2.Send("")
			},
			want: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"Set-Cookie":   {"name=agus; Path=/"},
			},
		},
		{
			name: "Expire",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus").Expires(expiration)
				r2.Cookie(c)
				c.GetExpires()
				c.GetRawExpires()
				r2.Send("")
			},
			want: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"Set-Cookie":   {"name=agus; Expires=" + expirationStr},
			},
		},
		{
			name: "Secure",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus").Secure(true)
				r2.Cookie(c)
				c.GetSecure()
				r2.Send("")
			},
			want: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"Set-Cookie":   {"name=agus; Secure"},
			},
		},
		{
			name: "SameSite",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Name("name").Value("agus").SameSite(http.SameSiteDefaultMode)
				r2.Cookie(c)
				c.GetSameSite()
				r2.Send("")
			},
			want: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"Set-Cookie":   {"name=agus"},
			},
		},
		{
			name: "Raw",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				c.Raw("name=agus")
				c.RawExpires("Wed, 23-Nov-2011 01:05:03 GMT")
				r2.Cookie(c)
				c.GetRaw()
				r2.Send("")
			},
			want: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
			},
		},
		{
			name: "Unparsed",
			handler: func(r1 Request, r2 Response) {
				c := Cookie{}
				unparsed := []string{"ok"}
				c.Unparsed(unparsed)
				r2.Cookie(c)
				c.GetUnparsed()
				r2.Send("")
			},
			want: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			tt.handler.ServeHTTP(w, req, nil, nil, nil, map[string]interface{}{})
			resp := w.Result()
			if header := resp.Header; !reflect.DeepEqual(header, tt.want) {
				t.Errorf("resp = %v, want %v", header, tt.want)
			}
		})
	}
}
