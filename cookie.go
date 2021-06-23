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

func (k *Cookie) Name(name string) *Cookie {
	k.c.Name = name
	return k
}

func (k *Cookie) Path(path string) *Cookie {
	k.c.Path = path
	return k
}

func (k *Cookie) Value(value string) *Cookie {
	k.c.Value = value
	return k
}

func (k *Cookie) Expires(time time.Time) *Cookie {
	k.c.Expires = time
	return k
}

func (k *Cookie) MaxAge(maxAge int) *Cookie {
	k.c.MaxAge = maxAge
	return k
}

func (k *Cookie) HttpOnly(httpOnly bool) *Cookie {
	k.c.HttpOnly = httpOnly
	return k
}

func (k *Cookie) Secure(secure bool) *Cookie {
	k.c.Secure = secure
	return k
}

func (k *Cookie) SameSite(sameSite http.SameSite) *Cookie {
	k.c.SameSite = sameSite
	return k
}

func (k *Cookie) Raw(raw string) *Cookie {
	k.c.Raw = raw
	return k
}

func (k *Cookie) RawExpires(expires string) *Cookie {
	k.c.RawExpires = expires
	return k
}

func (k *Cookie) Unparsed(Unparsed []string) *Cookie {
	k.c.Unparsed = Unparsed
	return k
}

func (k *Cookie) cookie() *http.Cookie {
	return &k.c
}
