package main

import (
	"fmt"
)

type Plugin1 struct{}

func (p *Plugin1) Run(args string) int {
	fmt.Printf("Plugin1 report: %v\n", args)
	return 0
}

func Create() interface{} {
	return new(Plugin1)
}
