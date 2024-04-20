//go:build !dev
// +build !dev

package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed build/browser
var _assets embed.FS

//go:embed build/server/index.wasm
var wasmBundle []byte

var _corrected, _ = fs.Sub(_assets, "build")

func Files(prefix string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.FS(_corrected)))
}
