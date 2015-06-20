package main

import (
	"github.com/mcartmell/tippler/tippler"
)

func main() {
	tippler.LoadAreas()
	go tippler.RunServer()
	tippler.RunStream()
}
