package main

import (
	"line-oa-backend/modules/servers"
)

func main() {
	server := servers.NewServer()
	server.Start()
}
