package main

import "embed-asset/web"

func main() {
	web.StartServer("0.0.0.0:12345")
}
