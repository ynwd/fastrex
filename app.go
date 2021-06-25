package fastrex

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type App interface {
	// Sets containers
	Container(map[string]interface{}) App
	// Routes HTTP GET requests to the specified path with the specified callback functions.
	Get(string, Handler, ...Middleware) App
	// Routes HTTP CONNECT requests to the specified path with the specified callback functions
	Connect(string, Handler, ...Middleware) App
	// Routes HTTP DELETE requests to the specified path with the specified callback functions
	Delete(string, Handler, ...Middleware) App
	// Routes HTTP HEAD requests to the specified path with the specified callback functions
	Head(string, Handler, ...Middleware) App
	// Routes HTTP PUT requests to the specified path with the specified callback functions
	Put(string, Handler, ...Middleware) App
	// Routes HTTP PATCH requests to the specified path with the specified callback functions
	Patch(string, Handler, ...Middleware) App
	// Routes HTTP TRACE requests to the specified path with the specified callback functions
	Trace(string, Handler, ...Middleware) App
	// Routes HTTP POST requests to the specified path with the specified callback functions
	Post(string, Handler, ...Middleware) App
	// Mounts the specified middleware function
	Use(Middleware) App
	// Sets static files
	Static(folder string, path ...string) App
	// Sets a logger
	Log(*log.Logger) App
	// Sets a context
	Ctx(context.Context) App
	// Binds and listens for connections on the specified host and port.
	Listen(port int, args ...interface{}) error
	// ListenServerless dispatches the request to the handler whose pattern most closely matches the request URL.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	// Sets app serverless
	Serverless(bool) App
	// Sets a host name
	Host(string) App
	// ParseFiles creates a new Template and parses the template definitions from the named files.
	Template(path string) App
	// SetKeepAlivesEnabled controls whether HTTP keep-alives are enabled. By default, keep-alives
	// are always enabled. Only very resource-constrained environments or servers in the process of
	// shutting down should disable them.
	SetKeepAlivesEnabled(v bool)
	// RegisterOnShutdown registers a function to call on Shutdown.
	// This can be used to gracefully shutdown connections that have undergone ALPN protocol upgrade or
	// that have been hijacked. This function should start protocol-specific graceful shutdown,
	// but should not wait for shutdown to complete.
	RegisterOnShutdown(f func())
	// Close immediately closes all active net.Listeners and any connections in state StateNew,
	// StateActive, or StateIdle. For a graceful shutdown, use Shutdown.
	//
	// Close does not attempt to close (and does not even know about) any hijacked connections,
	// such as WebSockets.
	//
	// Close returns any error returned from closing the Server's underlying Listener(s).
	Close() error
	// Shutdown gracefully shuts down the server without interrupting any active connections.
	// Shutdown works by first closing all open listeners, then closing all idle connections,
	// and then waiting indefinitely for connections to return to idle and then shut down.
	// If the provided context expires before the shutdown is complete, Shutdown returns
	// the context's error, otherwise it returns any error returned from closing the Server's
	// underlying Listener(s).
	//
	// When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS immediately return
	// ErrServerClosed. Make sure the program doesn't exit and waits instead for Shutdown to return.
	//
	// Shutdown does not attempt to close nor wait for hijacked connections such as WebSockets.
	// The caller of Shutdown should separately notify such long-lived connections of shutdown and
	// wait for them to close, if desired. See RegisterOnShutdown for a way to register shutdown
	// notification functions.
	//
	// Once Shutdown has been called on a server, it may not be reused; future calls to methods
	// such as Serve will return ErrServerClosed.
	Shutdown(ctx context.Context)
}

type Handler func(Request, Response)

type routerKey string

type Middleware func(Request, Response, Next)

type Next func(Request, Response)

type ErrMiddleware struct {
	Error error
	Code  int
}

type appRoute struct {
	path        string
	method      string
	handler     Handler
	middlewares []Middleware
}

type app struct {
	logger       *log.Logger
	server       *http.Server
	template     *template.Template
	ctx          context.Context
	container    map[string]interface{}
	routes       map[string]appRoute
	filename     []string
	middlewares  []Middleware
	serverless   bool
	staticFolder string
	staticPath   string
	host         string
}

const (
	errMiddlewareKey = routerKey("error")
)

// Create App instance
func New() App {
	return &app{
		container:    make(map[string]interface{}),
		routes:       map[string]appRoute{},
		middlewares:  []Middleware{},
		server:       &http.Server{},
		staticFolder: "",
		staticPath:   "",
	}
}

func (r *app) Container(c map[string]interface{}) App {
	if c != nil {
		r.container = c
	}
	return r
}

func (r *app) Use(m Middleware) App {
	if m != nil {
		r.middlewares = append(r.middlewares, m)
	}
	return r
}

func (r *app) Log(logger *log.Logger) App {
	r.logger = logger
	return r
}

func (r *app) Ctx(ctx context.Context) App {
	r.ctx = ctx
	return r
}

func (r *app) handler(serverless bool) http.Handler {
	return &httpHandler{r.routes, r.middlewares, r.template, r.logger, r.ctx, r.staticFolder, r.staticPath, serverless}
}

func (r *app) Close() error {
	return r.server.Close()
}

func (r *app) RegisterOnShutdown(f func()) {
	r.server.RegisterOnShutdown(f)
}

func (r *app) SetKeepAlivesEnabled(v bool) {
	r.server.SetKeepAlivesEnabled(v)
}

func (r *app) Shutdown(ctx context.Context) {
	r.server.Shutdown(ctx)
}

func (r *app) listenAndServe(addr string) error {
	r.server = &http.Server{
		Addr:    addr,
		Handler: r.handler(false),
	}

	return r.server.ListenAndServe()
}

func (r *app) listenAndServeTLS(addr string, certFile string, keyFile string) error {
	r.server = &http.Server{
		Addr:    addr,
		Handler: r.handler(false),
	}
	return r.server.ListenAndServeTLS(certFile, keyFile)
}

func (r *app) Serverless(serverless bool) App {
	r.serverless = serverless
	return r
}

func (r *app) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if len(r.filename) > 0 {
		r.handleTemplate()
	}
	h := r.handler(true)
	h.ServeHTTP(res, req)
}

func (r *app) Listen(port int, args ...interface{}) error {
	if len(r.filename) > 0 {
		err := r.handleTemplate()
		if err != nil {
			return err
		}
	}
	if len(args) == 0 || len(args) == 1 {
		return r.handleNonTLS(port, args)
	} else if len(args) == 2 || len(args) == 3 {
		return r.handleTLS(port, args)
	} else {
		return errors.New("error: invalid args")
	}
}

func (r *app) handleTemplate() error {
	if r.serverless {
		r.transform()
	}

	tmpl, err := template.ParseFiles(r.filename...)
	if err != nil {
		return err
	}
	r.template = tmpl
	return nil
}

func (r *app) transform() {
	filename := make([]string, 0)
	for _, v := range r.filename {
		filename = append(filename, serverlessFolder+v)
	}
	r.filename = filename
}

func (r *app) handleNonTLS(port int, args []interface{}) error {
	var (
		addr     string
		callback func(err error)
		ok       bool
	)

	if r.host == "" {
		addr = "localhost:" + strconv.Itoa(port)
	} else {
		addr = r.host + ":" + strconv.Itoa(port)
	}

	if len(args) == 1 {
		callback, ok = args[0].(func(e error))
		if !ok {
			return errors.New("error: invalid callback")
		}
		callback(nil)
	} else {
		fmt.Printf("Listening on http://%v \n", addr)
	}

	err := r.listenAndServe(addr)
	if err != nil {
		if callback != nil {
			callback(err)
		} else {
			fmt.Println(err)
		}
	}

	return err
}

func (r *app) handleTLS(port int, args []interface{}) error {
	var (
		addr     string
		ok       bool
		callback func(err error)
	)
	if r.host == "" {
		addr = "localhost:" + strconv.Itoa(port)
	} else {
		addr = r.host + ":" + strconv.Itoa(port)
	}
	var certFile, keyFile string
	cert, ok := args[0].(string)
	if ok {
		certFile = cert
	} else {
		return errors.New("error: invalid certFile")
	}
	key, ok := args[1].(string)
	if ok {
		keyFile = key
	} else {
		return errors.New("error: invalid keyFile")
	}

	if len(args) == 3 {
		callback, ok = args[2].(func(error))
		if !ok {
			return errors.New("error: invalid callback")
		}
		callback(nil)
	} else {
		fmt.Printf("Listening on https://%v \n", addr)
	}

	err := r.listenAndServeTLS(addr, certFile, keyFile)
	if err != nil {
		if callback != nil {
			callback(err)
		} else {
			fmt.Println(err)
		}
	}

	return err
}

func (r *app) Template(filename string) App {
	r.filename = append(r.filename, filename)
	return r
}

func (r *app) Host(host string) App {
	r.host = host
	return r
}

func (r *app) Static(folder string, path ...string) App {
	r.staticFolder = folder
	length := len(path)
	if length == 0 {
		r.staticPath = slash
	} else if length == 1 {
		r.staticPath = path[0]
	}
	return r
}

func appendMiddleware(m []Middleware) []Middleware {
	if len(m) == 0 {
		return nil
	}
	mid := []Middleware{}
	mid = append(mid, m...)
	return mid
}

func (r *app) Get(path string, handler Handler, middleware ...Middleware) App {
	route := appRoute{path, http.MethodGet, handler, appendMiddleware(middleware)}
	key := http.MethodGet + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Connect(path string, handler Handler, middleware ...Middleware) App {
	route := appRoute{path, http.MethodConnect, handler, appendMiddleware(middleware)}
	key := http.MethodConnect + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Delete(path string, handler Handler, middleware ...Middleware) App {
	route := appRoute{path, http.MethodDelete, handler, appendMiddleware(middleware)}
	key := http.MethodDelete + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Head(path string, handler Handler, middleware ...Middleware) App {
	route := appRoute{path, http.MethodHead, handler, appendMiddleware(middleware)}
	key := http.MethodHead + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Put(path string, handler Handler, middleware ...Middleware) App {
	route := appRoute{path, http.MethodPut, handler, appendMiddleware(middleware)}
	key := http.MethodPut + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Patch(path string, handler Handler, middleware ...Middleware) App {
	route := appRoute{path, http.MethodPatch, handler, appendMiddleware(middleware)}
	key := http.MethodPatch + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Trace(path string, handler Handler, middleware ...Middleware) App {
	route := appRoute{path, http.MethodTrace, handler, appendMiddleware(middleware)}
	key := http.MethodTrace + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Post(path string, handler Handler, middleware ...Middleware) App {
	route := appRoute{path, http.MethodPost, handler, appendMiddleware(middleware)}
	key := http.MethodPost + splitter + path
	r.routes[key] = route
	return r
}
