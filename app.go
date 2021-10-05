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

// Fastrex ..
type Fastrex func(app App) App

// App instance
type App interface {
	// Add app module
	Register(app Fastrex, url ...string) App
	// Sets Dependency
	Add(name string, i interface{}) App
	// Get Dependency
	Dependency(name string) interface{}
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

	Routes() map[string]AppRoute

	Middleware() []Middleware

	Templates() []string

	StaticFolder() string

	StaticPath() string
}

// Handler ...
type Handler func(Request, Response)

type routerKey string

// Middleware ...
type Middleware func(Request, Response, Next)

// Next ...
type Next func(Request, Response)

// ErrMiddleware ...
type ErrMiddleware struct {
	Error error
	Code  int
}

// AppRoute ...
type AppRoute struct {
	path        string
	method      string
	handler     Handler
	middlewares []Middleware
}

type app struct {
	logger     *log.Logger
	server     *http.Server
	template   *template.Template
	ctx        context.Context
	container  map[string]interface{}
	routes     map[string]AppRoute
	filename   []string
	serverless bool

	host              string
	apps              map[string]App
	middlewares       []Middleware
	moduleMiddlewares map[string][]Middleware

	staticFolder       string
	moduleStaticFolder map[string]string

	staticPath       string
	moduleStaticPath map[string]string
}

const (
	errMiddlewareKey = routerKey("error")
)

// New ...
func New() App {
	return &app{
		apps:               map[string]App{},
		container:          map[string]interface{}{},
		routes:             map[string]AppRoute{},
		middlewares:        []Middleware{},
		moduleMiddlewares:  map[string][]Middleware{},
		server:             &http.Server{},
		staticFolder:       "",
		moduleStaticFolder: map[string]string{},
		staticPath:         "",
		moduleStaticPath:   map[string]string{},
	}
}

func (r *app) StaticFolder() string {
	return r.staticFolder
}

func (r *app) StaticPath() string {
	return r.staticPath
}

func (r *app) Templates() []string {
	return r.filename
}

func (r *app) Middleware() []Middleware {
	return r.middlewares
}

func (r *app) Routes() map[string]AppRoute {
	return r.routes
}

func (r *app) Register(app Fastrex, url ...string) App {
	newApp := app(New())
	if len(url) > 0 {
		r.apps[url[0]] = newApp
		return r
	}

	r.apps[""] = newApp
	return r
}

func (r *app) Add(name string, i interface{}) App {
	r.container[name] = i
	return r
}

func (r *app) Dependency(name string) interface{} {
	return r.container[name]
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

func (r *app) mutate() {
	for url, app := range r.apps {
		newPath := url
		if len(app.Middleware()) > 0 {
			r.moduleMiddlewares[newPath] = app.Middleware()
		}
		if len(app.StaticFolder()) > 0 {
			r.moduleStaticFolder[newPath] = app.StaticFolder()
		}
		if len(app.StaticPath()) > 0 {
			r.moduleStaticPath[newPath] = app.StaticPath()
		}
		r.filename = append(r.filename, app.Templates()...)
		for _, route := range app.Routes() {
			if route.path == "/" {
				route.path = ""
			}
			newPath = url + route.path
			newKey := route.method + splitter + newPath
			newRoute := AppRoute{
				path:        newPath,
				method:      route.method,
				handler:     route.handler,
				middlewares: route.middlewares,
			}
			r.routes[newKey] = newRoute
			if len(app.Middleware()) > 0 {
				r.moduleMiddlewares[newPath] = app.Middleware()
			}
		}

	}

	if len(r.filename) > 0 {
		err := r.handleTemplate()
		if err != nil {
			panic(err)
		}
	}
}

func (r *app) handler(serverless bool) http.Handler {
	if len(r.apps) > 0 {
		r.mutate()
	}

	return &httpHandler{
		apps:               r.apps,
		container:          r.container,
		routes:             r.routes,
		template:           r.template,
		logger:             r.logger,
		ctx:                r.ctx,
		staticFolder:       r.staticFolder,
		moduleStaticFolder: r.moduleStaticFolder,
		staticPath:         r.staticPath,
		moduleStaticPath:   r.moduleStaticPath,
		serverless:         serverless,
		middlewares:        r.middlewares,
		moduleMiddlewares:  r.moduleMiddlewares,
	}
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
	err := r.server.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
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
		err := r.handleTemplate()
		if err != nil && r.logger != nil {
			r.logger.Println(err)
		}
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
	route := AppRoute{path, http.MethodGet, handler, appendMiddleware(middleware)}
	key := http.MethodGet + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Connect(path string, handler Handler, middleware ...Middleware) App {
	route := AppRoute{path, http.MethodConnect, handler, appendMiddleware(middleware)}
	key := http.MethodConnect + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Delete(path string, handler Handler, middleware ...Middleware) App {
	route := AppRoute{path, http.MethodDelete, handler, appendMiddleware(middleware)}
	key := http.MethodDelete + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Head(path string, handler Handler, middleware ...Middleware) App {
	route := AppRoute{path, http.MethodHead, handler, appendMiddleware(middleware)}
	key := http.MethodHead + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Put(path string, handler Handler, middleware ...Middleware) App {
	route := AppRoute{path, http.MethodPut, handler, appendMiddleware(middleware)}
	key := http.MethodPut + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Patch(path string, handler Handler, middleware ...Middleware) App {
	route := AppRoute{path, http.MethodPatch, handler, appendMiddleware(middleware)}
	key := http.MethodPatch + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Trace(path string, handler Handler, middleware ...Middleware) App {
	route := AppRoute{path, http.MethodTrace, handler, appendMiddleware(middleware)}
	key := http.MethodTrace + splitter + path
	r.routes[key] = route
	return r
}

func (r *app) Post(path string, handler Handler, middleware ...Middleware) App {
	route := AppRoute{path, http.MethodPost, handler, appendMiddleware(middleware)}
	key := http.MethodPost + splitter + path
	r.routes[key] = route
	return r
}
