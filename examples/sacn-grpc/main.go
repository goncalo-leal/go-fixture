package main

import (
	"fmt"

	"github.com/goncalo-leal/go-fixture"
)

func main() {
	f := fixture.LoadFromFile("config.json")
	fmt.Println(f.Start())
}
