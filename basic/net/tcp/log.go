package tcp

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "DEBUG ", log.LstdFlags|log.Lshortfile|log.Ltime)
)
