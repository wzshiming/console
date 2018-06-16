package router

import (
	"github.com/wzshiming/console"
	docker "github.com/wzshiming/console/docker"
	shell "github.com/wzshiming/console/shell"
	ssh "github.com/wzshiming/console/ssh"
)

func init() {
	console.Register("shell", shell.NewShellSessions)
	console.Register("docker", docker.NewDockerSessions)
	console.Register("ssh", ssh.NewSshSessions)
}
