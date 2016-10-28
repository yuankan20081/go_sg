// gate project main.go
package main

import (
	"fmt"
	"io"
	"os"
	"session"
	"time"
)

func main() {
	f, err := os.Open("./main.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	s := session.New(f, 0)
	go s.Flush()
	io.Copy(os.Stdout, s)
	s.Stop()
	time.Sleep(time.Second * 1)
}
