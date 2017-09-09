package session_docker

import "github.com/wzshiming/console"

func init() {
	console.Register("docker", NewDockerSessions)
}
