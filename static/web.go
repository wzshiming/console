package static

import (
	"embed"
)

//go:embed css
//go:embed js
//go:embed index.html
//go:embed robot.txt
var Web embed.FS
