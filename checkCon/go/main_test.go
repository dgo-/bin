package main

import "testing"

func Test_CheckInput(t *testing.T) {
    _, e := CheckInput("www.google.com:443")
    if e != nil {
      t.Fail()
    }
}

