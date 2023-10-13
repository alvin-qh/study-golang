package main

import (
	"fmt"

	"alvin.study/plugin/proto"
)

type Plugin1 struct{}

func (p *Plugin1) Run(args string) int {
	fmt.Printf("Plugin1 report: %v\n", args)
	return 0
}

func Create() proto.Runnable {
	return new(Plugin1)
}
