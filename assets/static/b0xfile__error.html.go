// Code generaTed by fileb0x at "2019-03-06 15:05:53.48097074 +0500 +05 m=+0.028043901" from config file "fileb0x.yml" DO NOT EDIT.
// modified(2018-05-26 13:29:04 +0500 +05)
// original path: assets/error.html

package static

import (
  
  "os"
)

// FileErrorHTML is "/error.html"
var FileErrorHTML = []byte("\x3c\x73\x65\x63\x74\x69\x6f\x6e\x20\x63\x6c\x61\x73\x73\x3d\x22\x73\x65\x63\x74\x69\x6f\x6e\x22\x3e\x0a\x20\x20\x20\x20\x3c\x64\x69\x76\x20\x63\x6c\x61\x73\x73\x3d\x22\x62\x6f\x78\x20\x68\x61\x73\x2d\x62\x61\x63\x6b\x67\x72\x6f\x75\x6e\x64\x2d\x77\x61\x72\x6e\x69\x6e\x67\x22\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x70\x3e\x7b\x65\x72\x72\x6f\x72\x7d\x3c\x2f\x70\x3e\x0a\x20\x20\x20\x20\x3c\x2f\x64\x69\x76\x3e\x0a\x3c\x2f\x73\x65\x63\x74\x69\x6f\x6e\x3e")

func init() {
  

  f, err := FS.OpenFile(CTX, "/error.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    panic(err)
  }

  
  _, err = f.Write(FileErrorHTML)
  if err != nil {
    panic(err)
  }
  

  err = f.Close()
  if err != nil {
    panic(err)
  }
}

