package main

import (
	"fmt"

	"alvin.study/plugin/proto"
)

type Plugin2 struct{}

func (p *Plugin2) Run(args string) int {
	fmt.Printf("Plugin2 report: %v\n", args)
	return 0
}

func Create() proto.Runnable {
	return new(Plugin2)
}
