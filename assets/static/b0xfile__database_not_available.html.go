// Code generaTed by fileb0x at "2019-10-12 15:06:05.13536 +0500 +05 m=+0.024291708" from config file "fileb0x.yml" DO NOT EDIT.
// modified(2019-04-16 16:30:46 +0500 +05)
// original path: assets/database_not_available.html

package static

import (
  
  "os"
)

// FileDatabaseNotAvailableHTML is "/database_not_available.html"
var FileDatabaseNotAvailableHTML = []byte("\x3c\x73\x65\x63\x74\x69\x6f\x6e\x20\x63\x6c\x61\x73\x73\x3d\x22\x73\x65\x63\x74\x69\x6f\x6e\x22\x3e\x0a\x20\x20\x20\x20\x3c\x64\x69\x76\x20\x63\x6c\x61\x73\x73\x3d\x22\x6e\x6f\x74\x69\x66\x69\x63\x61\x74\x69\x6f\x6e\x20\x69\x73\x2d\x64\x61\x6e\x67\x65\x72\x22\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x68\x33\x3e\x3c\x73\x74\x72\x6f\x6e\x67\x3e\x44\x61\x74\x61\x62\x61\x73\x65\x20\x6e\x6f\x74\x20\x61\x76\x61\x69\x6c\x61\x62\x6c\x65\x3c\x2f\x73\x74\x72\x6f\x6e\x67\x3e\x3c\x2f\x68\x33\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x70\x3e\x53\x6f\x6d\x65\x74\x68\x69\x6e\x67\x20\x77\x65\x6e\x74\x20\x77\x72\x6f\x6e\x67\x20\x77\x68\x69\x6c\x65\x20\x74\x72\x79\x69\x6e\x67\x20\x74\x6f\x20\x63\x6f\x6e\x6e\x65\x63\x74\x20\x74\x6f\x20\x64\x61\x74\x61\x62\x61\x73\x65\x2e\x20\x43\x68\x65\x63\x6b\x20\x6c\x6f\x67\x73\x20\x66\x6f\x72\x20\x64\x65\x74\x61\x69\x6c\x73\x2e\x3c\x2f\x70\x3e\x0a\x20\x20\x20\x20\x3c\x2f\x64\x69\x76\x3e\x0a\x3c\x2f\x73\x65\x63\x74\x69\x6f\x6e\x3e")

func init() {
  

  f, err := FS.OpenFile(CTX, "/database_not_available.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    panic(err)
  }

  
  _, err = f.Write(FileDatabaseNotAvailableHTML)
  if err != nil {
    panic(err)
  }
  

  err = f.Close()
  if err != nil {
    panic(err)
  }
}

