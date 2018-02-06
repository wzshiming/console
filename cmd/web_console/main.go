package main

import (
	"flag"

	"github.com/urfave/negroni"
	static "github.com/wzshiming/console/cmd/web_console/bind_static"
)

//go:generate go-bindata -o bind_static/static_binddata.go -pkg static static/...

var port = flag.String("addr", "0.0.0.0:8888", "Listen for server")

func main() {
	flag.Parse()

	// web
	n := negroni.New(
		negroni.NewLogger(),
		negroni.NewRecovery(),
	)
	// n.Use(negroni.NewStatic(http.Dir("static")))
	n.Use(negroni.NewStatic(static.NewFileSystem()))

	router, err := execRouter()
	if err != nil {
		panic(err)
	}

	n.UseHandler(router)
	n.Run(*port)
}
