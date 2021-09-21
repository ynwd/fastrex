package fastrex_test

import (
	"github.com/fastrodev/fastrex"
)

func ExampleListen(port int) {
	helloHandler := func(req fastrex.Request, res fastrex.Response) {
		res.Send("Hello world")
	}
	app := fastrex.New()
	app.Get("/", helloHandler)
	err := app.Listen(9000)
	if err != nil {
		panic(err)
	}
}

func ExampleListen_callback() {
	helloHandler := func(req fastrex.Request, res fastrex.Response) {
		res.Send("Hello world")
	}
	app := fastrex.New()
	app.Get("/", helloHandler)
	err := app.Listen(9000, func(error) {
		print("Listening on %v\n", 9000)
	})
	if err != nil {
		panic(err)
	}
}

func ExampleListen_TLS() {
	helloHandler := func(req fastrex.Request, res fastrex.Response) {
		res.Send("Hello world")
	}
	app := fastrex.New()
	app.Get("/", helloHandler)
	err := app.Listen(9000, "cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}
}

func ExampleListen_TLS_callback() {
	helloHandler := func(req fastrex.Request, res fastrex.Response) {
		res.Send("Hello world")
	}
	app := fastrex.New()
	app.Get("/", helloHandler)
	err := app.Listen(9000, "cert.pem", "key.pem", func(error) {
		print("Listening on %v\n", 9000)
	})
	if err != nil {
		panic(err)
	}
}

func ExampleTemplate() {
	app := fastrex.New()
	app.Template("index.html")
	err := app.Listen(9000)
	if err != nil {
		panic(err)
	}
}
