package main

import (
	"flag"

	"github.com/dimeko/wserver/client"
	"github.com/dimeko/wserver/server"
)

func main() {
	usage := flag.String("usage", "client", "choose client or server")
	host := flag.String("host", "localhost", "the host to initiate connection with")
	port := flag.String("port", "1337", "the port of the ws connection")
	path := flag.String("path", "ws", "the path for the ws connection")
	name := flag.String("name", "client", "name of the client")
	flag.Parse()

	if *usage == "client" {
		client.StartClient(*host, *port, *path, *name)
	} else {
		server.StartServer()
	}
}
