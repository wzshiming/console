package shell

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"unsafe"

	"github.com/wzshiming/console"
)

type SessionsShell struct {
	sessions map[string]*exec.Cmd
}

var _ console.Sessions = (*SessionsShell)(nil)

func NewShellSessions(host string) (console.Sessions, error) {
	return &SessionsShell{
		sessions: map[string]*exec.Cmd{},
	}, nil
}

func (d *SessionsShell) CreateExec(req *console.ReqCreateExec) (*console.RespCreateExec, error) {
	cli := exec.Command(req.Cmd)
	id := "0x" + strconv.FormatUint(uint64(uintptr(unsafe.Pointer(cli))), 16)
	d.sessions[id] = cli
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
	}()

	cli.Stdin = io.TeeReader(NewReader(ws), ws)
	cli.Stdout = ws
	cli.Stderr = ws

	return cli.Run()
}

func (d *SessionsShell) ResizeExecTTY(req *console.ReqResizeExecTTY) error {
	return nil
}

type Reader struct {
	buf io.Reader
}

func NewReader(read io.Reader) *Reader {
	return &Reader{read}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.buf.Read(p)
	if err != nil {
		return 0, err
	}
	buf := p[:n]
	buf = bytes.Replace(buf, []byte{13}, []byte{'\r', '\n'}, -1)
	copy(p[:len(buf)], buf)
	return len(buf), nil
}
