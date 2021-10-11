package main

import (
	"github-network/database/pkg/server"
)

func main() {
	s := server.InitServer(8080)
	s.Run()
}
