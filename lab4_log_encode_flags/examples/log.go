package main

import (
	"bytes"
	"fmt"
	"log"
)

func main() {
	buf := &bytes.Buffer{}

	logger := log.New(buf, "demo time: ", log.LstdFlags)
	logger.Println("hello, world")
	logger.Println("bye")

	fmt.Printf("%s\n", buf)
}
