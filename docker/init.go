package docker

import (
	"github.com/wzshiming/console"
)

func init() {
	console.Register("docker", NewDockerSessions)
}
