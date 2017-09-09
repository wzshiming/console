package session_docker

import (
	"io"

	"github.com/fsouza/go-dockerclient"
	"github.com/wzshiming/console"
	"github.com/wzshiming/ffmt"
)

type SessionsDocker struct {
	cli *docker.Client
}

var _ = (*SessionsDocker)(nil)

func NewDockerSessions(host string) (console.Sessions, error) {

	cli, err := docker.NewClient(host)
	if err != nil {
		return nil, err
	}

	return &SessionsDocker{
		cli: cli,
	}, nil
}

func (d *SessionsDocker) CreateExec(req *console.ReqCreateExec) (*console.RespCreateExec, error) {
	// 创建连接
	exec, err := d.cli.CreateExec(docker.CreateExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Env:          nil,
		Cmd:          []string{req.Cmd},
		Container:    req.CId,
		User:         "",
		Privileged:   true,
	})

	if err != nil {
		return nil, err
	}

	return &console.RespCreateExec{
		EId: exec.ID,
	}, nil
}

func (d *SessionsDocker) StartExec(id string, ws io.ReadWriter) error {
	// 执行连接
	err := d.cli.StartExec(id, docker.StartExecOptions{
		InputStream:  ws,
		OutputStream: ws,
		ErrorStream:  ws,
		Detach:       false,
		Tty:          true,
		RawTerminal:  true,
	})

	if err != nil {
		ffmt.Mark(err)
		return err
	}

	return nil
}

func (d *SessionsDocker) ResizeExecTTY(req *console.ReqResizeExecTTY) error {
	err := d.cli.ResizeExecTTY(req.EId, req.Height, req.Width)
	if err != nil {
		ffmt.Mark(err)
		return err
	}

	return nil
}
