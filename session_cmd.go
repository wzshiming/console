package console

import (
	"fmt"
	"io"
	"os/exec"
	"unsafe"
)

type SessionsCmd struct {
	sessions map[string]*exec.Cmd
}

var _ = (*SessionsCmd)(nil)

func NewCmdSessions(host string) (Sessions, error) {
	return &SessionsCmd{
		sessions: map[string]*exec.Cmd{},
	}, nil
}

func (d *SessionsCmd) CreateExec(req *ReqCreateExec) (*RespCreateExec, error) {
	cli := exec.Command(req.Cmd)
	id := fmt.Sprint(unsafe.Pointer(cli))
	d.sessions[id] = cli
	return &RespCreateExec{
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

func (d *SessionsCmd) ResizeExecTTY(req *ReqResizeExecTTY) error {

	return nil
}
