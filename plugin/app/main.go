package main

import (
	"fmt"
	"plugin"
)

type Runnable interface {
	Run(args string) int
}

func main() {
	runner := loadPlugin("./p1.so")
	runner.Run("Hello")

	runner = loadPlugin("./p2.so")
	runner.Run("OK")
}

func loadPlugin(plugName string) Runnable {
	plug, err := plugin.Open(plugName)
	if err != nil {
		panic(err)
	}

	symb, err := plug.Lookup("Create")
	if err != nil {
		panic(err)
	}

	factory, ok := symb.(func() any)
	if !ok {
		panic(fmt.Errorf("symbol `Create` cannot cast to function"))
	}
	return factory().(Runnable)
}
