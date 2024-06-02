package tcp

import (
	"log"
	"os"
)

var (
	sLog = log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile|log.Ltime)
	cLog = log.New(os.Stdout, "[CLIENT] ", log.LstdFlags|log.Lshortfile|log.Ltime)
)
