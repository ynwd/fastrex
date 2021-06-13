package fastrex

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	slash    = "/"
	splitter = ":"
	notFound = "!"
	empty    = ""
)

type httpHandler struct {
	routes       map[string]httpRoute
	middlewares  []Middleware
	template     *template.Template
	logger       *log.Logger
	ctx          context.Context
	staticFolder string
	staticPath   string
	serverless   bool
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
	}

	if h.ctx != nil {
		r = r.WithContext(h.ctx)
	}

	key := h.getRouteKey(r.Method, r.URL.Path)
	route, ok := h.routes[key]
	if !ok {
		folder := h.staticFolder
		if h.serverless {
			folder = serverlessFolder + h.staticFolder
		}
		http.StripPrefix(h.staticPath, http.FileServer(http.Dir(folder))).ServeHTTP(w, r)
		return
	}

	if len(h.middlewares) > 0 || len(route.middlewares) > 0 {
		h.handleMiddleware(route, w, r)
	} else if route.handler != nil {
		route.handler(*newRequest(r), newResponse(w, r, h.template))
	}
}

func (h *httpHandler) handleMiddleware(route httpRoute,
	w http.ResponseWriter, r *http.Request) {
	var (
		next     bool
		request  Request
		response Response
	)
	lengthOfRouteMiddleware := len(route.middlewares)
	lengthOfAppMiddleware := len(h.middlewares)

	if lengthOfAppMiddleware > 0 {
		next, request, response = h.loopMiddleware(route, h.middlewares, w, r, lengthOfAppMiddleware)
		if !next {
			return
		}
	}

	if lengthOfRouteMiddleware > 0 {
		next, request, response = h.loopMiddleware(route, route.middlewares, w, r, lengthOfRouteMiddleware)
	}

	if next {
		route.handler(request, response)
	}
}

func (h *httpHandler) loopMiddleware(
	route httpRoute,
	middlewares []Middleware,
	w http.ResponseWriter, r *http.Request,
	length int) (bool, Request, Response) {
	var (
		next     bool
		request  Request
		response Response
	)
	for i := range middlewares {
		responseMid := newResponse(w, r, h.template)
		requestMid := newRequest(r)
		middlewares[length-1-i](
			*requestMid,
			responseMid,
			func(req Request, res Response) {
				next = true
				request = req
				response = res
				e, ok := req.Context().Value(errMiddlewareKey).(ErrMiddleware)
				if ok {
					next = false
					http.Error(w, e.Error.Error(), e.Code)
				}
			})
	}
	return next, request, response
}

func (h *httpHandler) getRouteKey(incomingMethod string, incomingPath string) string {
	for _, r := range h.routes {
		if incomingMethod == r.method && h.validate(r.path, incomingPath) {
			return r.method + splitter + r.path
		}
	}
	return notFound
}

func (h *httpHandler) validate(path string, incoming string) bool {
	p := h.split(path)
	i := h.split(incoming)
	if len(p) != len(i) {
		return false
	}
	return parsePath(p, i)
}

func (h *httpHandler) split(str string) []string {
	var s []string
	s = strings.Split(str, "/")
	s = append(s[:0], s[1:]...)
	return s
}

func isAllTrue(a []bool) bool {
	for i := 0; i < len(a); i++ {
		if !a[i] {
			return false
		}
	}
	return true
}

func parsePath(p []string, incoming []string) (valid bool) {
	var results []bool
	for idx, path := range p {
		results = append(results, isValidPath(idx, path, incoming))
	}
	valid = isAllTrue(results)
	return
}

func isValidPath(idx int, path string, incoming []string) bool {
	if incoming[idx] == path || regex(incoming[idx], path) {
		return true
	}
	return false
}

func regex(incoming string, path string) bool {
	if (incoming != empty) && strings.Contains(path, splitter) {
		if strings.Contains(path, "(") && strings.Contains(path, ")") {
			r := regexp.MustCompile(getPattern(path))
			return r.MatchString(incoming)
		}
		return true
	}
	return false
}

func getPattern(s string) (str string) {
	i := strings.Index(s, "(")
	if i >= 0 {
		j := strings.Index(s, ")")
		if j >= 0 {
			str = s[i+1 : j]
		}
	}
	return
}

type HandlerFunc func(Request, Response)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(*newRequest(r), newResponse(w, r, nil))
}