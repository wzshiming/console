package ssh

import (
	"github.com/wzshiming/console"
)

func init() {
	console.Register("ssh", NewSshSessions)
}
