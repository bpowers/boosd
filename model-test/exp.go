package main

import (
	"boosd/runtime"
	"log"
)

func main() {
	var err error
	if err = runtime.Init(); err != nil {
		log.Fatalf("runtime.Init(): %s", err)
	}
}