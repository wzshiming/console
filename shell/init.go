package shell

import (
	"github.com/wzshiming/console"
	_ "github.com/wzshiming/winseq"
)

func init() {
	console.Register("shell", NewShellSessions)
}
