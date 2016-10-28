// gate project main.go
package main

import (
	"fmt"
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
	s.OnReadFunc(func(p []byte, sid int) {
		fmt.Println(string(p))
	})
	s.Serve(time.Second * 1000)
}
