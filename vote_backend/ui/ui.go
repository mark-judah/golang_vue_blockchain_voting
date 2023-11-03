package resource

import "embed"

//go:embed public/index.html
var Html []byte

//go:emmbed public/favicon.ico
var Favicon []byte

// go:embed public/static/css
var CssStatic embed.FS

// go:embed public/static/img
var ImgStatic embed.FS

// go:embed public/static/js
var JsStatic embed.FS

// go:embed public/static/fonts
var FontStatic embed.FS

//go:embed all:public
var Static embed.FS
