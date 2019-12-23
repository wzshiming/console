package session_shell

import (
	"github.com/wzshiming/console"
)

func init() {
	console.Register("shell",  NewShellSessions)
}

