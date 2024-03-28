package docker

const (
	DumpDockerfile = `
	FROM scratch
	COPY . .
	`
)
