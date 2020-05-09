package main

import "github.com/lucasmls/backend-cacautime/application/server"

func main() {
	s := server.NewService(server.ServiceInput{})

	s.Run()
}
