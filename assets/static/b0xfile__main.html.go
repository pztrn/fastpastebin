// Code generaTed by fileb0x at "2019-10-12 15:06:05.134184 +0500 +05 m=+0.023115110" from config file "fileb0x.yml" DO NOT EDIT.
// modified(2019-10-12 15:03:51.676112966 +0500 +05)
// original path: assets/main.html

package static

import (
  
  "os"
)

// FileMainHTML is "/main.html"
var FileMainHTML = []byte("\x3c\x21\x44\x4f\x43\x54\x59\x50\x45\x20\x68\x74\x6d\x6c\x3e\x0a\x3c\x68\x74\x6d\x6c\x3e\x0a\x0a\x3c\x68\x65\x61\x64\x3e\x0a\x20\x20\x20\x20\x3c\x6d\x65\x74\x61\x20\x63\x68\x61\x72\x73\x65\x74\x3d\x22\x75\x74\x66\x2d\x38\x22\x3e\x0a\x20\x20\x20\x20\x3c\x6d\x65\x74\x61\x20\x6e\x61\x6d\x65\x3d\x22\x76\x69\x65\x77\x70\x6f\x72\x74\x22\x20\x63\x6f\x6e\x74\x65\x6e\x74\x3d\x22\x77\x69\x64\x74\x68\x3d\x64\x65\x76\x69\x63\x65\x2d\x77\x69\x64\x74\x68\x2c\x20\x69\x6e\x69\x74\x69\x61\x6c\x2d\x73\x63\x61\x6c\x65\x3d\x31\x22\x3e\x0a\x20\x20\x20\x20\x3c\x74\x69\x74\x6c\x65\x3e\x46\x61\x73\x74\x20\x50\x61\x73\x74\x65\x20\x42\x69\x6e\x3c\x2f\x74\x69\x74\x6c\x65\x3e\x0a\x20\x20\x20\x20\x3c\x6c\x69\x6e\x6b\x20\x72\x65\x6c\x3d\x22\x73\x74\x79\x6c\x65\x73\x68\x65\x65\x74\x22\x20\x68\x72\x65\x66\x3d\x22\x2f\x73\x74\x61\x74\x69\x63\x2f\x63\x73\x73\x2f\x62\x75\x6c\x6d\x61\x2d\x30\x2e\x37\x2e\x35\x2e\x6d\x69\x6e\x2e\x63\x73\x73\x22\x3e\x0a\x20\x20\x20\x20\x3c\x6c\x69\x6e\x6b\x20\x72\x65\x6c\x3d\x22\x73\x74\x79\x6c\x65\x73\x68\x65\x65\x74\x22\x20\x68\x72\x65\x66\x3d\x22\x2f\x73\x74\x61\x74\x69\x63\x2f\x63\x73\x73\x2f\x62\x75\x6c\x6d\x61\x2d\x74\x6f\x6f\x6c\x74\x69\x70\x2d\x33\x2e\x30\x2e\x30\x2e\x6d\x69\x6e\x2e\x63\x73\x73\x22\x3e\x0a\x20\x20\x20\x20\x3c\x6c\x69\x6e\x6b\x20\x72\x65\x6c\x3d\x22\x73\x74\x79\x6c\x65\x73\x68\x65\x65\x74\x22\x20\x68\x72\x65\x66\x3d\x22\x2f\x73\x74\x61\x74\x69\x63\x2f\x63\x73\x73\x2f\x73\x74\x79\x6c\x65\x2e\x63\x73\x73\x22\x3e\x0a\x3c\x2f\x68\x65\x61\x64\x3e\x0a\x0a\x3c\x62\x6f\x64\x79\x3e\x0a\x20\x20\x20\x20\x7b\x6e\x61\x76\x69\x67\x61\x74\x69\x6f\x6e\x7d\x20\x7b\x64\x6f\x63\x75\x6d\x65\x6e\x74\x42\x6f\x64\x79\x7d\x0a\x3c\x2f\x62\x6f\x64\x79\x3e\x0a\x3c\x66\x6f\x6f\x74\x65\x72\x20\x63\x6c\x61\x73\x73\x3d\x22\x66\x6f\x6f\x74\x65\x72\x22\x3e\x0a\x20\x20\x20\x20\x7b\x66\x6f\x6f\x74\x65\x72\x7d\x0a\x3c\x2f\x66\x6f\x6f\x74\x65\x72\x3e\x0a\x0a\x3c\x2f\x68\x74\x6d\x6c\x3e")

func init() {
  

  f, err := FS.OpenFile(CTX, "/main.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    panic(err)
  }

  
  _, err = f.Write(FileMainHTML)
  if err != nil {
    panic(err)
  }
  

  err = f.Close()
  if err != nil {
    panic(err)
  }
}

