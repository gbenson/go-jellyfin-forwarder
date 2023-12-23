package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go-jellyfin-forwarder HOST")
		os.Exit(2)
	}

	log.SetFlags(log.Lshortfile)
	err := Run(os.Args[1])
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func Run(server string) error {
	const want = "who is JellyfinServer?"
	const port = "7359"
	conn, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		return err
	}
	//defer conn.Close() // XXX maybe?

	server += ":" + port
	log.Println("listening")
	for {
		var buf [1024]byte
		got := buf[:]
		n, client, err := conn.ReadFrom(got)
		log.Println("got", n, client, err, got)
		if n == len(want) && string(got[:n]) == want {
			log.Println("hello:", client, server)
			//go Forward(
			continue
		}
		if err != nil {
			log.Println("warning:", err)
			continue
		}
	}
}
