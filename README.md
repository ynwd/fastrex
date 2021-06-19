# Fastrex
[![][build]](https://github.com/fastrodev/fastrex/actions/workflows/build.yml) [![Coverage Status][cov]](https://coveralls.io/github/fastrodev/fastrex?branch=main) [![][reference]](https://pkg.go.dev/github.com/fastrodev/fastrex?tab=doc)

Fast and simple web application framework for Go inspired by the most popular node.js web framework: Express.js. It implements `ServeHTTP` interface so you can use express style routing. It also wraps and extends the net/http `Request` and `ResponseWriter` into an easy to read and use function signature. 

## Get Started
Init folder and install:
```
mkdir app && cd app
go mod init github.com/fastrodev/examples
go get github.com/fastrodev/fastrex
```
Create main.go file:
```go
package main

import "github.com/fastrodev/fastrex"

func handler(req fastrex.Request, res fastrex.Response) {
  res.Json(`{"message":"hello"}`)
}

func createApp() fastrex.App {
  app := fastrex.New()
  app.Get("/", handler)
  return app
}

func main() {
  createApp().Listen(9000)
}

```

Run webapp locally:
```
go run main.go
```
## Serverless deployment

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
  createApp().ServeHTTP(w, r)
}

```
How to deploy:
```
gcloud functions deploy Main --runtime go113 --trigger-http --allow-unauthenticated
```
Demo and full example: [https://github.com/fastrodev/serverless](https://github.com/fastrodev/serverless)
## Contributing
We appreciate your help! The main purpose of this repository is to improve performance and readability, making it faster and easier to use.

[build]: https://github.com/fastrodev/fastrex/actions/workflows/build.yml/badge.svg
[reference]: https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white "reference"
[cov]: https://coveralls.io/repos/github/fastrodev/fastrex/badge.svg?branch=main

