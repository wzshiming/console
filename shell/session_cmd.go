package session_shell

import (
	"fmt"
	"io"
	"os/exec"
	"unsafe"

	"github.com/wzshiming/console"
)

type SessionsCmd struct {
	sessions map[string]*exec.Cmd
}

var _ = (*SessionsCmd)(nil)

func NewCmdSessions(host string) (console.Sessions, error) {
	return &SessionsCmd{
		sessions: map[string]*exec.Cmd{},
	}, nil
}

func (d *SessionsCmd) CreateExec(req *console.ReqCreateExec) (*console.RespCreateExec, error) {
	cli := exec.Command(req.Cmd)
	id := fmt.Sprint(unsafe.Pointer(cli))
	d.sessions[id] = cli
	return &console.RespCreateExec{
		EId: id,
	}, nil
}

func (d *SessionsCmd) StartExec(eid string, ws io.ReadWriter) error {
	cli, ok := d.sessions[eid]
	if !ok {
		return fmt.Errorf("Can not find eid " + eid)
	}
	defer func() {
		delete(d.sessions, eid)
	}()

	cli.Stdin = ws
	cli.Stdout = ws
	cli.Stderr = ws

	return cli.Run()
}

func (d *SessionsCmd) ResizeExecTTY(req *console.ReqResizeExecTTY) error {

	return nil
}
