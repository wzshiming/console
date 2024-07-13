//go:build !windows
// +build !windows

package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"unsafe"

	"github.com/creack/pty"
	"github.com/wzshiming/console"
)

type SessionsShell struct {
	sessions map[string]*os.File
}

var _ console.Sessions = (*SessionsShell)(nil)

func NewShellSessions(host string) (console.Sessions, error) {
	return &SessionsShell{
		sessions: map[string]*os.File{},
	}, nil
}

func (d *SessionsShell) CreateExec(req *console.ReqCreateExec) (*console.RespCreateExec, error) {
	sh := exec.Command(req.Cmd)
	id := "0x" + strconv.FormatUint(uint64(uintptr(unsafe.Pointer(sh))), 16)
	// Start the command with a pty.
	ptmx, err := pty.Start(sh)
	if err != nil {
		return nil, err
	}

	d.sessions[id] = ptmx
	return &console.RespCreateExec{
		EId: id,
	}, nil
}

func (d *SessionsShell) StartExec(eid string, ws io.ReadWriter) error {
	cli, ok := d.sessions[eid]
	if !ok {
		return fmt.Errorf("Can not find eid " + eid)
	}
	defer func() {
		delete(d.sessions, eid)
		cli.Close()
	}()

	go io.Copy(cli, ws)
	io.Copy(ws, cli)
	return nil
}

func (d *SessionsShell) ResizeExecTTY(req *console.ReqResizeExecTTY) error {
	cli, ok := d.sessions[req.EId]
	if !ok {
		return fmt.Errorf("Can not find eid " + req.EId)
	}
	return pty.Setsize(cli, &pty.Winsize{
		Rows: uint16(req.Height),
		Cols: uint16(req.Width),
	})
}
