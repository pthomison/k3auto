package hack

import (
	"embed"
)

var (
	//go:embed *.yaml
	K3dConfig embed.FS
)
