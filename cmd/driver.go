package session_cmd

import "github.com/wzshiming/console"

func init() {
	console.Register("cmd", NewCmdSessions)
}
