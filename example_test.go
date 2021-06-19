package fastrex_test

import (
	"github.com/fastrodev/fastrex"
)

func ExampleListen() {
	helloHandler := func(req fastrex.Request, res fastrex.Response) {
		res.Send("Hello world")
	}
	app := fastrex.New()
	app.Get("/", helloHandler)
	app.Listen(9000)
}

func ExampleListen_callback() {
	helloHandler := func(req fastrex.Request, res fastrex.Response) {
		res.Send("Hello world")
	}
	app := fastrex.New()
	app.Get("/", helloHandler)
	app.Listen(9000, func(error) {
		print("Listening on %v\n", 9000)
	})
}

func ExampleListen_TLS() {
	helloHandler := func(req fastrex.Request, res fastrex.Response) {
		res.Send("Hello world")
	}
	app := fastrex.New()
	app.Get("/", helloHandler)
	app.Listen(9000, "cert.pem", "key.pem")
}

func ExampleListen_TLS_callback() {
	helloHandler := func(req fastrex.Request, res fastrex.Response) {
		res.Send("Hello world")
	}
	app := fastrex.New()
	app.Get("/", helloHandler)
	app.Listen(9000, "cert.pem", "key.pem", func(error) {
		print("Listening on %v\n", 9000)
	})
}

func ExampleTemplate() {
	app := fastrex.New()
	app.Template("index.html")
	app.Listen(9000)
}
