package web

import "embed"

/*
  Below code comments are likely not just a normal comments.
  Dont mess with them unless you know what youre doing.
  https://pkg.go.dev/embed
*/

//go:embed all:dist
var EmbeddedFiles embed.FS
