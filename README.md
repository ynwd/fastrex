# Fastrex
[![][build]](https://github.com/fastrodev/fastrex/actions/workflows/build.yml) [![Coverage Status][cov]](https://coveralls.io/github/fastrodev/fastrex?branch=main) [![][reference]](https://pkg.go.dev/github.com/fastrodev/fastrex?tab=doc)

Fast and simple web application framework for Go inspired by the most popular node.js web framework: Express.js. It implements `ServeHTTP` interface so you can use express style routing. It also wraps and extends the net/http `Request` and `ResponseWriter` into an easy to read and use function signature. 

* [Get started](#get-started)
* [Middleware](#middleware)
* [Module](#module)
* [Template](#template)
* [Serverless](#serverless)
* [Benchmarks](#benchmarks)

## Get Started
Init folder and install:
```
mkdir app && cd app
go mod init app
go get github.com/fastrodev/fastrex
```
Create main.go file:
```go
package main

import "github.com/fastrodev/fastrex"

func handler(req fastrex.Request, res fastrex.Response) {
	res.Send("root")
}

func main() {
	app := fastrex.New()
	app.Get("/", handler)
	err := app.Listen(9000)
	if err != nil {
		panic(err)
	}
}


```

Run webapp locally:
```
go run main.go
```

## Middleware
You can access `Request` and `Response` field and function before the handler process the incoming request.
### App Middleware
```go
package main

import "github.com/fastrodev/fastrex"

func handler(req fastrex.Request, res fastrex.Response) {
	res.Send("root")
}

func appMiddleware(req fastrex.Request, res fastrex.Response, next fastrex.Next) {
	if req.URL.Path == "/" {
		res.Send("appMiddleware")
		return
	}

	next(req, res)
}

func main() {
	app := fastrex.New()
	app.Use(appMiddleware)
	app.Get("/", handler)
	err := app.Listen(9000)
	if err != nil {
		panic(err)
	}
}

```

### Route Middleware

```go
package main

import "github.com/fastrodev/fastrex"

func handler(req fastrex.Request, res fastrex.Response) {
	res.Send("root")
}

func routeMiddleware(req fastrex.Request, res fastrex.Response, next fastrex.Next) {
	if req.URL.Path == "/" {
		res.Send("appMiddleware")
		return
	}
}

func main() {
	app := fastrex.New()
	app.Get("/", handler, routeMiddleware)
	err := app.Listen(9000)
	if err != nil {
		panic(err)
	}
}

```

## Module
You can group urls, routes, and handlers with `Register()`.
```go
package main

import "github.com/fastrodev/fastrex"

func module(app fastrex.App) fastrex.App {
	moduleMiddleware := func(req fastrex.Request, res fastrex.Response, next fastrex.Next) {
		if req.URL.Path == "/api/user" {
			res.Send("userMiddleware")
			return
		}
		next(req, res)
	}
	handler := func(req fastrex.Request, res fastrex.Response) {
		res.Send("userModule")
	}
	app.Use(moduleMiddleware)
	app.Get("/user", handler)
	return app
}

func main() {
	app := fastrex.New()
	app.Register(module, "/api")
	err := app.Listen(9000)
	if err != nil {
		panic(err)
	}
}

```
## Template
You can render html by create HTML template at `template` folder.
```html
<html>{{.Title}}{{.Name}}</html>
```
Then you add them to app with `Template` function. 

And finally, call `Render` function from handler.
```go
package main

import "github.com/fastrodev/fastrex"

func handler(req fastrex.Request, res fastrex.Response) {
	data := struct {
		Title string
		Name  string
	}{
		"hallo",
		"world",
	}
	err := res.Render(data)
	if err != nil {
		panic(err)
	}
}

func main() {
	app := fastrex.New()
	app.Template("template/app.html")
	app.Get("/", handler)
	err := app.Listen(9000)
	if err != nil {
		panic(err)
	}
}

```
## Serverless

You can deploy your codes to [google cloud function](https://cloud.google.com/functions). With this approach, you don't call the `Listen` function again. You must create a new function as the entry point for standard net/http `Request` and` ResponseWriter`.

```go
package serverless

import (
  "net/http"

  "github.com/fastrodev/fastrex"
)

func handler(req fastrex.Request, res fastrex.Response) {
  res.Json(`{"message":"hello"}`)
}

func createApp() fastrex.App {
  app := fastrex.New()
  app.Get("/", handler)
  return app
}

func Main(w http.ResponseWriter, r *http.Request) {
  createApp().Serverless(true).ServeHTTP(w, r)
}

```
How to deploy:
```
gcloud functions deploy Main --runtime go113 --trigger-http --allow-unauthenticated
```
Demo and full example: [`https://github.com/fastrodev/serverless`](https://github.com/fastrodev/serverless)

## Benchmarks
|Module|Requests/sec|Transfer/sec|
|--|--:|--:|
|Fastrex|95249.11|10.99MB|
|Go std|88700.49|10.83MB|
|Node std|50696.05|6.48MB|
|Express|9006.68|2.05MB|

Benchmarks repository: [`https://github.com/fastrodev/benchmarks`](https://github.com/fastrodev/benchmarks)

## Contributing
We appreciate your help! The main purpose of this repository is to improve performance and readability, making it faster and easier to use.

[build]: https://github.com/fastrodev/fastrex/actions/workflows/build.yml/badge.svg
[reference]: https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white "reference"
[cov]: https://coveralls.io/repos/github/fastrodev/fastrex/badge.svg?branch=main

