package main

import (
  "io"
  "os"
  "strings"
)

type rot13Reader struct {
  r io.Reader
}

func (rot *rot13Reader) Read(b []byte) (numByte int, err error) {
  numByte, err = rot.r.Read(b)
  for i := range b {
    if b[i] >= 'A' && b[i] < 'Z' {
      b[i] = 'A' + ((b[i] - 'A' + 13) % 26)
    }
    if b[i] >= 'a' && b[i] <= 'z' {
      b[i] = 'a' + ((b[i] - 'a' + 13) % 26)
    }
  }
  return
}

func main() {
  s := strings.NewReader("Lbh penpxrq gur pbqr!")
  r := rot13Reader{s}
  io.Copy(os.Stdout, &r)
}