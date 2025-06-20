package main

import (
	"fmt"
)

type Plugin2 struct{}

func (p *Plugin2) Run(args string) int {
	fmt.Printf("Plugin2 report: %v\n", args)
	return 0
}

func Create() any {
	return new(Plugin2)
}
