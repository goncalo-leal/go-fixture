package main

import (
	"fmt"
	"time"

	"github.com/goncalo-leal/go-fixture"
)

func main() {
	f := fixture.LoadFromFile("config.json")
	fmt.Println(f.Start())
	for {
		time.Sleep(1 * time.Second)
	}
}
