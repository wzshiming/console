package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/urfave/negroni"
	"github.com/wzshiming/console"
	"github.com/wzshiming/console/router"
	"github.com/wzshiming/console/static"
)

var port = flag.Int("p", 8888, "Listen port")
var ip = flag.String("ip", "0.0.0.0", "Listen ip")
var name = flag.String("name", "shell", "Driver name, shell docker and ssh")
var host = flag.String("host", "", "Connect to the host address, It has to be docker or SSH Driver")
var cid = flag.String("cid", "", "Docker Container id, It has to be docker Driver")
var cmd = flag.String("cmd",
	func() string {
		if runtime.GOOS != "windows" {
			return "sh"
		}
		return "cmd"
	}(),
	"command to execute")
var disable = flag.Bool("d", false, "Disable url parameters")

func main() {
	flag.Parse()

	// web
	n := negroni.New(
		negroni.NewLogger(),
		negroni.NewRecovery(),
	)

	n.Use(negroni.NewStatic(static.NewFileSystem()))
	n.UseHandler(router.ExecRouter(*disable, &console.ReqCreateExec{
		Name: *name,
		Host: *host,
		CId:  *cid,
		Cmd:  *cmd,
	}))
	n.Run(fmt.Sprintf("%v:%v", *ip, *port))
}
