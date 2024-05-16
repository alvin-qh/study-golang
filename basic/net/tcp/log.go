package tcp

import (
	"log"
	"os"
)

var (
	lServer = log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile|log.Ltime)
	lClient = log.New(os.Stdout, "[CLIENT] ", log.LstdFlags|log.Lshortfile|log.Ltime)
)
