package fastrex

import (
	"net/http"
	"time"
)

type Cookie struct {
	c http.Cookie
}

func (k *Cookie) Domain(name string) *Cookie {
	k.c.Domain = name
	return k
}

func (k *Cookie) GetDomain() string {
	return k.c.Domain
}

func (k *Cookie) Name(name string) *Cookie {
	k.c.Name = name
	return k
}

func (k *Cookie) GetName() string {
	return k.c.Name
}

func (k *Cookie) Path(path string) *Cookie {
	k.c.Path = path
	return k
}

func (k *Cookie) GetPath() string {
	return k.c.Path
}

func (k *Cookie) Value(value string) *Cookie {
	k.c.Value = value
	return k
}

func (k *Cookie) GetValue() string {
	return k.c.Value
}

func (k *Cookie) Expires(time time.Time) *Cookie {
	k.c.Expires = time
	return k
}

func (k *Cookie) GetExpires() time.Time {
	return k.c.Expires
}

func (k *Cookie) MaxAge(maxAge int) *Cookie {
	k.c.MaxAge = maxAge
	return k
}

func (k *Cookie) GetMaxAge() int {
	return k.c.MaxAge
}

func (k *Cookie) HttpOnly(httpOnly bool) *Cookie {
	k.c.HttpOnly = httpOnly
	return k
}

func (k *Cookie) GetHttpOnly() bool {
	return k.c.HttpOnly
}

func (k *Cookie) Secure(secure bool) *Cookie {
	k.c.Secure = secure
	return k
}

func (k *Cookie) GetSecure() bool {
	return k.c.Secure
}

func (k *Cookie) SameSite(sameSite http.SameSite) *Cookie {
	k.c.SameSite = sameSite
	return k
}

func (k *Cookie) GetSameSite() http.SameSite {
	return k.c.SameSite
}

func (k *Cookie) Raw(raw string) *Cookie {
	k.c.Raw = raw
	return k
}

func (k *Cookie) GetRaw() string {
	return k.c.Raw
}

func (k *Cookie) RawExpires(expires string) *Cookie {
	k.c.RawExpires = expires
	return k
}

func (k *Cookie) GetRawExpires() string {
	return k.c.RawExpires
}

func (k *Cookie) Unparsed(Unparsed []string) *Cookie {
	k.c.Unparsed = Unparsed
	return k
}

func (k *Cookie) GetUnparsed() []string {
	return k.c.Unparsed
}

func (k *Cookie) cookie() *http.Cookie {
	return &k.c
}
