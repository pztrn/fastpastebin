// Code generaTed by fileb0x at "2018-04-30 12:45:53.117322185 +0500 +05 m=+0.017381722" from config file "fileb0x.yml" DO NOT EDIT.
// modified(2018-04-30 12:45:50.434017895 +0500 +05)
// original path: assets/css/style.css

package static

import (
  
  "os"
)

// FileStaticCSSStyleCSS is "static/css/style.css"
var FileStaticCSSStyleCSS = []byte("\x23\x70\x61\x73\x74\x65\x2d\x63\x6f\x6e\x74\x65\x6e\x74\x73\x20\x7b\x0a\x20\x20\x20\x20\x68\x65\x69\x67\x68\x74\x3a\x20\x39\x30\x76\x68\x3b\x0a\x7d")

func init() {
  

  f, err := FS.OpenFile(CTX, "static/css/style.css", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    panic(err)
  }

  
  _, err = f.Write(FileStaticCSSStyleCSS)
  if err != nil {
    panic(err)
  }
  

  err = f.Close()
  if err != nil {
    panic(err)
  }
}

