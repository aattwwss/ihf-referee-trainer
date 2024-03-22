package public

import (
	"embed"
	"io/fs"
)

var (
	//go:embed static
	static embed.FS

	//go:embed html
	html embed.FS
)

func HTML() (fs.FS, error) {
	return fs.Sub(html, "html")
}

func Static() (fs.FS, error) {
	return fs.Sub(static, "static")
}
