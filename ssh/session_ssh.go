package ssh

import (
	"fmt"
	"io"
	"net/url"
	"strconv"
	"unsafe"

	"github.com/wzshiming/console"
	"golang.org/x/crypto/ssh"
)

type SessionsSsh struct {
	cli      *ssh.Client
	sessions map[string]*ssh.Session
}

var _ = (*SessionsSsh)(nil)

func NewSshSessions(host string) (console.Sessions, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	pwd := ""
	user := ""
	if u.User != nil {
		user = u.User.Username()
		pwd, _ = u.User.Password()
	}
	cli, err := ssh.Dial("tcp", u.Host, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		return nil, err
	}

	return &SessionsSsh{
		cli:      cli,
		sessions: map[string]*ssh.Session{},
	}, nil
}

func (d *SessionsSsh) CreateExec(req *console.ReqCreateExec) (*console.RespCreateExec, error) {
	cli, err := d.cli.NewSession()
	if err != nil {
		return nil, err
	}

	id := "0x" + strconv.FormatUint(uint64(uintptr(unsafe.Pointer(cli))), 16)
	d.sessions[id] = cli
	return &console.RespCreateExec{
		EId: id,
	}, nil
}

func (d *SessionsSsh) StartExec(eid string, ws io.ReadWriter) error {
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

	// Request pseudo terminal
	err := cli.RequestPty("xterm", 40, 80, ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	})
	if err != nil {
		return err
	}

	err = cli.Shell()
	if err != nil {
		return err
	}

	return cli.Wait()
}

func (d *SessionsSsh) ResizeExecTTY(req *console.ReqResizeExecTTY) error {
	cli, ok := d.sessions[req.EId]
	if !ok {
		return fmt.Errorf("Can not find eid " + req.EId)
	}
	return cli.WindowChange(req.Height, req.Width)
}
