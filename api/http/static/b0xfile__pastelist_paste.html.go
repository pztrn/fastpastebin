// Code generaTed by fileb0x at "2018-04-30 18:14:14.353534072 +0500 +05 m=+0.034414808" from config file "fileb0x.yml" DO NOT EDIT.
// modified(2018-04-30 18:14:11.780202908 +0500 +05)
// original path: assets/pastelist_paste.html

package static

import (
  
  "os"
)

// FilePastelistPasteHTML is "/pastelist_paste.html"
var FilePastelistPasteHTML = []byte("\x3c\x64\x69\x76\x20\x63\x6c\x61\x73\x73\x3d\x22\x63\x6f\x6e\x74\x65\x6e\x74\x22\x3e\x0a\x20\x20\x20\x20\x3c\x64\x69\x76\x20\x63\x6c\x61\x73\x73\x3d\x22\x63\x61\x72\x64\x22\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x68\x65\x61\x64\x65\x72\x20\x63\x6c\x61\x73\x73\x3d\x22\x63\x61\x72\x64\x2d\x68\x65\x61\x64\x65\x72\x22\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x70\x20\x63\x6c\x61\x73\x73\x3d\x22\x63\x61\x72\x64\x2d\x68\x65\x61\x64\x65\x72\x2d\x74\x69\x74\x6c\x65\x22\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x50\x61\x73\x74\x65\x20\x23\x7b\x70\x61\x73\x74\x65\x49\x44\x7d\x2c\x20\x70\x6f\x73\x74\x65\x64\x20\x6f\x6e\x20\x7b\x70\x61\x73\x74\x65\x44\x61\x74\x65\x7d\x20\x61\x6e\x64\x20\x74\x69\x74\x6c\x65\x64\x20\x61\x73\x20\x22\x7b\x70\x61\x73\x74\x65\x54\x69\x74\x6c\x65\x7d\x22\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x2f\x70\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x2f\x68\x65\x61\x64\x65\x72\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x64\x69\x76\x20\x63\x6c\x61\x73\x73\x3d\x22\x63\x61\x72\x64\x2d\x63\x6f\x6e\x74\x65\x6e\x74\x22\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x64\x69\x76\x20\x63\x6c\x61\x73\x73\x3d\x22\x63\x6f\x6e\x74\x65\x6e\x74\x22\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x70\x72\x65\x3e\x7b\x70\x61\x73\x74\x65\x44\x61\x74\x61\x7d\x3c\x2f\x70\x72\x65\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x2f\x64\x69\x76\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x2f\x64\x69\x76\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x66\x6f\x6f\x74\x65\x72\x20\x63\x6c\x61\x73\x73\x3d\x22\x63\x61\x72\x64\x2d\x66\x6f\x6f\x74\x65\x72\x22\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x61\x20\x68\x72\x65\x66\x3d\x22\x2f\x70\x61\x73\x74\x65\x2f\x7b\x70\x61\x73\x74\x65\x49\x44\x7d\x22\x20\x63\x6c\x61\x73\x73\x3d\x22\x63\x61\x72\x64\x2d\x66\x6f\x6f\x74\x65\x72\x2d\x69\x74\x65\x6d\x20\x62\x75\x74\x74\x6f\x6e\x20\x69\x73\x2d\x73\x75\x63\x63\x65\x73\x73\x20\x69\x73\x2d\x72\x61\x64\x69\x75\x73\x6c\x65\x73\x73\x22\x3e\x56\x69\x65\x77\x3c\x2f\x61\x3e\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x3c\x2f\x66\x6f\x6f\x74\x65\x72\x3e\x0a\x20\x20\x20\x20\x3c\x2f\x64\x69\x76\x3e\x0a\x3c\x2f\x64\x69\x76\x3e")

func init() {
  

  f, err := FS.OpenFile(CTX, "/pastelist_paste.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    panic(err)
  }

  
  _, err = f.Write(FilePastelistPasteHTML)
  if err != nil {
    panic(err)
  }
  

  err = f.Close()
  if err != nil {
    panic(err)
  }
}

