package main

import (
	"flag"
	"net/http"

	"github.com/urfave/negroni"
)

var port = flag.String("addr", "0.0.0.0:8888", "Listen for server")

func main() {
	flag.Parse()

	// web
	n := negroni.New(
		negroni.NewLogger(),
		negroni.NewRecovery(),
		negroni.NewStatic(http.Dir("static")),
	)

	router, err := execRouter()
	if err != nil {
		panic(err)
	}

	n.UseHandler(router)
	n.Run(*port)
}
