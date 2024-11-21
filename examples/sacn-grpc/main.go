package main

import (
	"github.com/goncalo-leal/go-fixture/fixture"
)

func main() {
	f := fixture.LoadFromFile("config.json")
	f.Start()
}
