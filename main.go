package main

import (
	"encoding/base64"
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/wzshiming/ffmt"
	"golang.org/x/net/websocket"

	docker "github.com/fsouza/go-dockerclient"
)

var port = flag.String("port", "8888", "Port for server")
var host = flag.String("host", "tcp://127.0.0.1:2735", "Docker host")

var contextPath = ""

func main() {
	flag.Parse()

	if cp := os.Getenv("CONTEXT_PATH"); cp != "" {
		contextPath = strings.TrimRight(cp, "/")
	}

	http.Handle(contextPath+"/exec/", websocket.Handler(ExecContainer))
	http.Handle(contextPath+"/", http.StripPrefix(contextPath+"/", http.FileServer(http.Dir("./"))))
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		panic(err)
	}
}

func ExecContainer(ws *websocket.Conn) {
	wsParams := strings.Split(ws.Request().URL.Path[len(contextPath+"/exec/"):], ",")
	container := wsParams[0]
	cmd, _ := base64.StdEncoding.DecodeString(wsParams[1])

	if container == "" {
		ws.Write([]byte("Container does not exist"))
		return
	}

	client, err := docker.NewClient(host)
	if err != nil {
		ffmt.Mark(err)
		return
	}

	exec, err := client.CreateExec(docker.CreateExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Env:          []string{},
		Cmd:          []string{string(cmd)},
		Container:    container,
		User:         "",
		Privileged:   true,
	})

	if err != nil {
		ffmt.Mark(err)
		return
	}

	eo := docker.StartExecOptions{
		InputStream:  ws,
		OutputStream: ws,
		ErrorStream:  ws,

		Detach:      false,
		Tty:         true,
		RawTerminal: true,
	}

	err = client.StartExec(exec.ID, eo)

	if err != nil {
		ffmt.Mark(err)
		return
	}
}
