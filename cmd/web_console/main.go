package main

import (
	"flag"
	"fmt"

	"github.com/urfave/negroni"
	"github.com/wzshiming/console/router"
	"github.com/wzshiming/console/static"

	_ "github.com/wzshiming/console/docker"
	_ "github.com/wzshiming/console/shell"
	_ "github.com/wzshiming/console/ssh"
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

	router, err := router.ExecRouter()
	if err != nil {
		panic(err)
	}

	n.UseHandler(router)
	n.Run(fmt.Sprintf("%v:%v", *ip, *port))
}
