// Code generaTed by fileb0x at "2018-05-01 00:26:24.973087795 +0500 +05 m=+0.032533011" from config file "fileb0x.yml" DO NOT EDIT.
// modified(2018-04-30 23:54:11 +0500 +05)
// original path: assets/pagination_link.html

package static

import (
  
  "os"
)

// FilePaginationLinkHTML is "/pagination_link.html"
var FilePaginationLinkHTML = []byte("\x3c\x6c\x69\x3e\x0a\x20\x20\x20\x20\x3c\x61\x20\x63\x6c\x61\x73\x73\x3d\x22\x70\x61\x67\x69\x6e\x61\x74\x69\x6f\x6e\x2d\x6c\x69\x6e\x6b\x22\x20\x61\x72\x69\x61\x2d\x6c\x61\x62\x65\x6c\x3d\x22\x47\x6f\x20\x74\x6f\x20\x70\x61\x67\x65\x20\x7b\x70\x61\x67\x65\x4e\x75\x6d\x7d\x22\x20\x68\x72\x65\x66\x3d\x22\x7b\x70\x61\x67\x69\x6e\x61\x74\x69\x6f\x6e\x4c\x69\x6e\x6b\x7d\x22\x3e\x7b\x70\x61\x67\x65\x4e\x75\x6d\x7d\x3c\x2f\x61\x3e\x0a\x3c\x2f\x6c\x69\x3e")

func init() {
  

  f, err := FS.OpenFile(CTX, "/pagination_link.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    panic(err)
  }

  
  _, err = f.Write(FilePaginationLinkHTML)
  if err != nil {
    panic(err)
  }
  

  err = f.Close()
  if err != nil {
    panic(err)
  }
}

