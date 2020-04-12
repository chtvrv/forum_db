package main

import (
	serv "github.com/chtvrv/forum_db/app/server"
)

func main() {
	server := new(serv.Server)

	server.Run()
}
