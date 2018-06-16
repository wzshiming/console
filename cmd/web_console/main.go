package main

import (
	"flag"
	"fmt"

	"github.com/urfave/negroni"
	"github.com/wzshiming/console/router"
	"github.com/wzshiming/console/static"
)

var port = flag.Int("p", 8888, "Listen port")
var ip = flag.String("ip", "0.0.0.0", "Listen ip")

func main() {
	flag.Parse()

	// web
	n := negroni.New(
		negroni.NewLogger(),
		negroni.NewRecovery(),
	)

	n.Use(negroni.NewStatic(static.NewFileSystem()))
	n.UseHandler(router.ExecRouter())
	n.Run(fmt.Sprintf("%v:%v", *ip, *port))
}
