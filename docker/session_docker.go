package session_docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/wzshiming/console"
)

type SessionsDocker struct {
	cli *client.Client
}

var _ = (*SessionsDocker)(nil)

func NewDockerSessions(host string) (console.Sessions, error) {
	cli, err := client.NewClient(host, "1.39", nil, nil)
	if err != nil {
		return nil, err
	}

	return &SessionsDocker{
		cli: cli,
	}, nil
}

func (d *SessionsDocker) CreateExec(req *console.ReqCreateExec) (*console.RespCreateExec, error) {
	// 创建连接

	exec, err := d.cli.ContainerExecCreate(context.Background(), req.CId, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{req.Cmd},
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
	hr, err := d.cli.ContainerExecAttach(context.Background(), id, types.ExecConfig{
		Detach: false,
		Tty:    true,
	})
	if err != nil {
		return err
	}
	defer hr.Close()
	go io.Copy(ws, hr.Conn)
	io.Copy(hr.Conn, ws)

	return nil
}

func (d *SessionsDocker) ResizeExecTTY(req *console.ReqResizeExecTTY) error {
	return d.cli.ContainerExecResize(context.Background(), req.EId, types.ResizeOptions{
		Height: uint(req.Height),
		Width:  uint(req.Width),
	})
	return nil
}
