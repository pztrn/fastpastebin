// Code generaTed by fileb0x at "2019-03-06 15:05:53.472324169 +0500 +05 m=+0.019397329" from config file "fileb0x.yml" DO NOT EDIT.
// modified(2018-05-01 17:43:40 +0500 +05)
// original path: assets/css/style.css

package static

import (
  
  "os"
)

// FileStaticCSSStyleCSS is "static/css/style.css"
var FileStaticCSSStyleCSS = []byte("\x23\x70\x61\x73\x74\x65\x2d\x63\x6f\x6e\x74\x65\x6e\x74\x73\x20\x7b\x0a\x20\x20\x20\x20\x66\x6f\x6e\x74\x2d\x66\x61\x6d\x69\x6c\x79\x3a\x20\x6d\x6f\x6e\x6f\x73\x70\x61\x63\x65\x3b\x0a\x20\x20\x20\x20\x66\x6f\x6e\x74\x2d\x73\x69\x7a\x65\x3a\x20\x30\x2e\x39\x72\x65\x6d\x3b\x0a\x20\x20\x20\x20\x68\x65\x69\x67\x68\x74\x3a\x20\x39\x30\x76\x68\x3b\x0a\x7d\x0a\x0a\x2e\x70\x61\x73\x74\x65\x2d\x64\x61\x74\x61\x20\x7b\x0a\x20\x20\x20\x20\x66\x6f\x6e\x74\x2d\x73\x69\x7a\x65\x3a\x20\x30\x2e\x39\x72\x65\x6d\x3b\x0a\x7d")

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

