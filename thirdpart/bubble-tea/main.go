package main

import (
	"os"
	"study/thirdpart/bubbletea/chat"
	"study/thirdpart/bubbletea/simple/counter"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := os.Args

	var model tea.Model = nil

	for _, arg := range args[1:] {
		switch arg {
		case "simple":
			model = counter.NewModel()
		case "chat":
			model = chat.NewModel()
		}
	}

	if model != nil {
		tea.NewProgram(model).Run()
	}
}
