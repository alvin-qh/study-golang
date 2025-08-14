package main

import (
	"fmt"
	"os"
	"plugin"
)

type Runnable interface {
	Run(args string) int
}

func makePluginPath() string {
	pluginPath := os.Getenv("PLUGIN_PATH")
	if pluginPath == "" {
		pluginPath = "../dist"
	}

	return pluginPath
}

func main() {
	pluginPath := makePluginPath()
	fmt.Println("PLUGIN_PATH:", pluginPath)

	runner := loadPlugin(fmt.Sprintf("%s/p1.so", pluginPath))
	runner.Run("Hello")

	runner = loadPlugin(fmt.Sprintf("%s/p2.so", pluginPath))
	runner.Run("OK")
}

func loadPlugin(path string) Runnable {
	plug, err := plugin.Open(path)
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
