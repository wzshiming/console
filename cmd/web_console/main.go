package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/wzshiming/console"
	"github.com/wzshiming/console/router"
	"github.com/wzshiming/console/static"

	_ "github.com/wzshiming/console/docker"
	_ "github.com/wzshiming/console/shell"
	_ "github.com/wzshiming/console/ssh"
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

func init() {
	flag.Parse()
}

func main() {
	// web
	mux0 := mux.NewRouter()

	mux0.PathPrefix("/api").Handler(http.StripPrefix("/api", router.ExecRouter(*disable, &console.ReqCreateExec{
		Name: *name,
		Host: *host,
		CId:  *cid,
		Cmd:  *cmd,
	})))

	mux0.PathPrefix("/").Handler(http.FileServer(static.NewFileSystem()))

	mux := handlers.RecoveryHandler()(mux0)
	mux = handlers.CombinedLoggingHandler(os.Stdout, mux)
	p := fmt.Sprintf("%v:%v", *ip, *port)
	fmt.Printf("Open http://%s/ with your browser.\n", p)
	err := http.ListenAndServe(p, mux)
	if err != nil {
		fmt.Println(err)
	}
}
